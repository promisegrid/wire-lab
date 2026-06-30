#![no_std]
#![no_main]

use core::panic::PanicInfo;

// Intent: Keep the first Renode slice to a minimal M4F firmware artifact so
// POC17 proves platform viability before radio or full protocol behavior.
// Source: DI-pokin
#[repr(C)]
pub struct VectorTable {
    initial_stack: u32,
    reset: extern "C" fn() -> !,
}

#[link_section = ".vector_table.reset_vector"]
#[no_mangle]
pub static RESET_VECTOR: VectorTable = VectorTable {
    initial_stack: 0x2003_0000,
    reset: Reset,
};

#[no_mangle]
pub extern "C" fn Reset() -> ! {
    loop {
        core::hint::spin_loop();
    }
}

#[panic_handler]
fn panic(_info: &PanicInfo) -> ! {
    loop {
        core::hint::spin_loop();
    }
}
