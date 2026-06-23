# DN-zaraz: Bintags LoRa Frame Budget for PromiseGrid Mesh

Status: research/design note requiring verification before normative use.

This note records POC17 and later constrained-radio design pressure from the
`bintags` Feather M4/RFM9x prior art. It is not a PromiseGrid protocol spec and
does not by itself lock regulatory, hardware, driver, packet, proof-size, MTU,
fragmentation, or store-and-forward decisions. Treat concrete values below as
implementation-planning inputs that must be verified against the target radio
driver, hardware, region, and simulator fidelity before they become normative.
Source: `DI-govat`.

- **Target Stack:** angelajt/bintags (Asset Tracking) +
  promisegrid/wire-lab (Decentralized Consensus)
- **Hardware Profile:** Adafruit M4 Express 3857 + Adafruit 3231 LoRa
  Radio FeatherWing (Semtech SX1276 / RFM95W)
- **Deployment Region:** United States (US915 Band, 902–928 MHz)
- **Topology Constraints:** 5-Hop Decentralized Mesh with
  Store-and-Forward Capabilities

---

## 1. Core Architectural Challenge & Mismatch

The goal is to transition `angelajt/bintags` from its native
point-to-point (P2P) star/broadcast topology into a decentralized,
consensus-driven mesh topology using `promisegrid/wire-lab`. 

While both systems utilize raw LoRa layers (avoiding heavy LoRaWAN
network overhead), they present a severe data footprint mismatch:

* **Bintags Baseline Payload:** Custom lean C/C++ structs map out to a
  tiny footprint of **11 to 23 bytes** per asset tracking or status
  update frame (Device ID, status flags, battery level, sequence
  counter, and sensor metadata).
* **PromiseGrid/Wire-Lab Overhead:** Operating as a consensus-based
  state machine, the wire-lab framework (specifically the
  secure-tokens and encrypted-payload proofs-of-concept) relies on
  cryptographic signatures (e.g., Ed25519 adding 64 bytes per block),
  state map dictionaries, and structured formatting like CBOR. This
  expands the base footprint to **128 to 256+ bytes**.
* **The Constraint:** The raw LoRa physical layer enforces strict
  Maximum Transmission Unit (MTU) limits dictated by hardware FIFO
  buffers and regional legal constraints.

---

## 2. LoRa Core Concepts & Constraints

### Spreading Factor (SF) Theory
Spreading Factor (SF) is the number of chirps used to encode a single
piece of data (a symbol). Values range from SF7 to SF12, with each
step up doubling the chirps per symbol.

* **Low SF (SF7):** High data rates, short Time-on-Air, low battery
  usage, but **shorter range**.
* **High SF (SF12):** Slow data rates, high Time-on-Air, high power
  usage, but **maximum range** due to increased receiver sensitivity.

### Regional Regulations & The US915 Dwell Time Limit

Under **FCC Part 15.247** regulations in the United States, standard
narrow-band channels (125 kHz) enforce a strict **400ms maximum dwell
time (Time-on-Air)** per packet transmission. This creates a hard
physical constraint: as your Spreading Factor goes up, the payload
size must drastically go down to stay legal.

| Spreading Factor / Bandwidth | Max Usable Payload Size | Legal / Regulatory Rule |
| :--- | :--- | :--- |
| **SF10 / 125 kHz** | ~11 Bytes | Limited by US 400ms Dwell Time |
| **SF9 / 125 kHz** | 53 Bytes | Limited by US 400ms Dwell Time |
| **SF7 / 125 kHz** | 242 Bytes | Fits within 400ms Dwell Time |
| **SF7 to SF9 / 500 kHz** | **242 Bytes** | **Exempt from Dwell Time Limit** via FCC Wideband Digital Modulation Mode |

*Note: 242 Bytes represents the hard absolute limit of the Semtech
SX1276's internal 255-byte FIFO buffer after accounting for mandatory
radio framing headers and CRC trailers.*

---

## 3. Optimizing for Maximum Payload Size via Store-and-Forward

Because PromiseGrid operates natively as an **asynchronous
store-and-forward architecture**, real-time transmission constraints
across the 5-hop mesh are bypassed. Data can be safely buffered
locally in flash memory and gossiped as the channel permits.

To get the absolute maximum payload size (**242 bytes**) without
incurring a 400ms legal airtime penalty, the mesh must utilize **FCC
Part 15.247 Wideband Digital Modulation Mode**. This requires locking
the radio into a **500 kHz bandwidth**.

### Comparative Performance Metrics (500 kHz Bandwidth)

| Metric | SF7 / 500 kHz Config | SF9 / 500 kHz Config |
| :--- | :--- | :--- |
| **Max Payload Space** | **242 Bytes** | **242 Bytes** |
| **Receiver Sensitivity** | ~ -111 dBm | **~ -117 dBm (+6 dB Link Budget Advantage)** |
| **Physical Range Factor** | Baseline (1x) | **~1.5x to 2x more range per hop** |
| **Time-on-Air (242B)** | ~35 milliseconds | ~142 milliseconds |
| **Battery Consumption** | Lowest | Medium-Low |
| **Mesh Collision Risk** | Minimal | Low |

### Operational Range Profiles

* **Deep Indoor / Warehouse:** 150m–400m (SF7) vs. 400m–1km (SF9).
  Heavily impacted by metal shelving and multipath reflections.
* **Dense Urban:** 500m–1.5km (SF7) vs. 1.5km–3km (SF9). Relies on the
  5-hop topology to route around building obstructions.
* **Line-of-Sight (Clear Outdoor):** 5km–12km (SF7) vs. 15+ km (SF9).

---

## 4. Hardware Driver & Code Implementation Layer

Since `bintags` already natively incorporates the **`Adafruit_RFM9x`**
driver library for its P2P operations, you do not need to replace your
low-level drivers. 

The `Adafruit_RFM9x` library wraps data frames inside an explicit
4-byte header sequence (`Destination Address`, `Sender Address`,
`Packet ID`, `Flags`) making it structurally identical to the popular
`RadioHead (RH_RF95)` P2P format. Because the driver manages these 4
bytes automatically, your clean application software buffer
(PromiseGrid state data + Bintags updates) must cap out at **238
bytes** to hit the exact 242-byte physical radio ceiling.

### Configuration Snippets

#### CircuitPython Realization
```python
import board
import busio
import digitalio
import adafruit_rfm9x

# Hardware pin bindings for Adafruit M4 Express 3857
spi = busio.SPI(board.SCK, board.MOSI, board.MISO)
cs = digitalio.DigitalInOut(board.RFM9X_CS)
reset = digitalio.DigitalInOut(board.RFM9X_RST)

# Initialize P2P Radio on US915 Core
rfm9x = adafruit_rfm9x.RFM9x(spi, cs, reset, 915.0)

# Apply Maximum Payload Wideband Mesh Settings
rfm9x.signal_bandwidth = 500000  # 500 kHz eliminates the FCC 400ms constraint
rfm9x.spreading_factor = 9       # High sensitivity for maximum structural range
rfm9x.coding_rate = 5            # Coding Rate 4/5 for low overhead

# Explicit Local Identity Mapping
rfm9x.node = 2                  # Local Mesh Node Identifier
rfm9x.destination = 3           # Explicit Next-Hop Routing Target
```

#### Arduino / C++ Realization
```cpp
#include <SPI.h>
#include <Adafruit_RFM9x.h>

#define RFM95_CS   8
#define RFM95_RST  4
#define RFM95_INT  3

Adafruit_RFM9x rfm9x(RFM95_CS, RFM95_INT);

void setup() {
    pinMode(RFM95_RST, OUTPUT);
    digitalWrite(RFM95_RST, HIGH);
    
    if (!rfm9x.init()) {
        while (1); // Hardware initialization fault trap
    }
    
    rfm9x.setFrequency(915.0);
    
    // Explicit register overrides to expand MTU buffer legally
    rfm9x.setSignalBandwidth(500000); // 500 kHz wideband
    rfm9x.setSpreadingFactor(9);      // SF9 depth profile
    rfm9x.setCodingRate4(5);          // CR 4/5 optimization
}
```

---

## 5. Overcoming the Mesh Limitations (Dynamic Options)

The Adafruit 3231 FeatherWing incorporates a **single physical modem
core (SX1276)**. Unlike multi-channel commercial gateways, it can only
demodulate exactly one Spreading Factor and one Bandwidth channel
configuration at a time. To implement a fluid topology over 5 hops,
three methods can be introduced into the firmware layer:

### A. Adaptive Data Rate (ADR) Link Fallback

Nodes default to fast transmissions via **SF7** to optimize battery
consumption and keep time-on-air down to **35ms**. If a targeted
next-hop node fails to reply with an Acknowledgment (ACK) header frame
within 3 sequential store-and-forward attempts, the sending M4 Express
dynamically alters its register (`rfm9x.spreading_factor = 9`) and
attempts remediation over the more robust **SF9** layer.

### B. Multi-SF Asynchronous Scanning

Idle nodes cycle their modem registers across an active polling loop
(e.g., listening on SF7 for 500ms, then executing an instruction to
hop to SF9 for 500ms). Transmitting mesh nodes compensate for this by
extending their packet's **preamble length** out past 1.0 seconds.
This guarantees that the receiver will naturally cycle into the
sender's exact SF layer while the introductory sync sequence is still
firing in the air.

### C. Content-Addressable "Pull" Synchronization

Instead of continuously pushing maximum-length 242-byte encrypted
structures along all 5 hops—which spikes channel congestion—nodes are
configured to only broadcast short, 32-byte cryptographic hashes
representing the current local root state of their PromiseGrid ledger
map. Neighboring nodes intercept this lean beacon frame; if they
detect an unrepresented ledger modification, they dynamically initiate
a localized P2P handoff command to explicitly request download of the
missing 238-byte block fragments.
