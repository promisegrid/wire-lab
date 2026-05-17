# Spec Requirement Section Scenarios

These scenarios are evidence for `TODO-kulih` 010.9 and `DR-robon`. They are not
a decision and not a final spec template. Source: `DI-pukap`.

## Scenario Matrix

| Scenario | Setup | What to test | Decision pressure |
|---|---|---|---|
| Minimal binding spec | Alice writes a low-level UDP-feed-like binding spec with a short normative surface. | Whether promise vocabulary and mental-model sections clarify or bloat a simple binding. | If even simple bindings benefit, required sections become stronger. |
| High-level promise protocol | Bob writes a group/session-like spec whose core semantics are human promises. | Whether layer-specific promise vocabulary prevents ambiguous implementation claims. | High-level protocols likely need explicit vocabulary even if lower layers do not. |
| 100-year review | Carol evaluates whether an old spec is still intelligible after tooling, hosts, and social assumptions changed. | Whether a required long-horizon section preserves the assumptions needed for future readers. | TE-dajot-style pressure may need a concrete spec-section home. |
| Layperson handoff | Dave, a non-kernel developer, decides whether he can trust an implementation based on the spec. | Whether layperson/easy-implementation summaries reduce misunderstanding without replacing normative text. | If the summary changes behavior expectations, it may need a normative hook. |
| Boilerplate failure | A spec includes all required sections but each is vague. | Whether review/conformance gates can reject empty prose. | Required sections need quality criteria or they become ceremony. |

## Expected Outputs

- Evidence for answering `DR-robon` without reopening the turn-177 replay.
- A recommended split between required spec sections, guide prose, and optional
  companion explanatory docs.
