# Blackbox expansion — execution state

Orchestrator memory. Updated after every step so a fresh session can resume.
Branch: `feat/blackbox-expansion` (from `feat/investigation-timeline`). Baseline at start: `go build ./...` OK; `go test ./internal/blackbox/... ./internal/server/...` OK (blackbox ~44s, server ~7min). Go 1.27.3 at `../.dev/toolchains/go/bin`.

Status legend: pending / dispatched / done / review / fixing / committed / blocked.
Reviewer: U = Ultimate, P = Performance. Tier: E = Efficient, P = Performance.

## Wave 1
| Packet | Status | Tier | Wave | Review | Commit |
|---|---|---|---|---|---|
| BBX-1-01 budget | committed | P | 1 | U: 1 HIGH livelock fixed + re-reviewed OK | 06add92 |
| BBX-1-02 engine | committed | P | 1 | U: clean (CounterReading rename accepted) | f11137c |
| BBX-1-03 config targets | committed | E | 1 | U: clean | cf63b04 |
| BBX-1-04 group types | committed | E | 1 | U: clean | f11137c |

## Wave 2
| BBX-1-05 group registry | committed | P | 2 | U: clean | (wave 2 merge) |
| BBX-1-06 PDH source | committed | P | 2 | U: 1 CRITICAL (PdhExpandWildCard 5th arg) fixed + arity re-verified | (wave 2 merge + fix) |
| BBX-1-10 wizard endpoint | committed | P | 2 | U: clean (pool run-state is live check) | (wave 2 merge) |

## Wave 3
| BBX-1-07 budget loop | committed | P | 3 | U: HIGH CPU-unit + MED false-gap fixed (r1); race + reader-late fixed (r2, atomic effective) | (wave 3 + fix r1/r2) |
| BBX-1-08 resolver proc | committed | P | 3 | U: LOW Start==0 cache fixed (r1) | (wave 3 + fix r1) |
| BBX-1-09 benchmark gate | committed | P | 3 | U: clean (0.335ms/tick, 86 allocs) | (wave 3) |
| BBX-1-12 wizard picker | committed | P | 3 | U: clean (1-line wizard mode prop + scn1 guards accepted) | (wave 3) |

## Wave 4
| BBX-1-11 resolver map | committed | P | 4 | U: HIGH all-or-nothing cache + HIGH stale PDH cache + wiring gap fixed; re-reviewed OK | (wave 4 + fix) |
| BBX-2-01 HTTP.sys | committed | P | 4 | U: clean | (wave 4) |
| BBX-2-03 disk+CPU | committed | P | 4 | U: clean; pre-existing clumsy test race stabilized (test-only) | (wave 4) |
| BBX-2-04 SQL host counters | committed | P | 4 | U: MED bRefresh 2nd call fixed; Windows enum wired; dedup pdhDetailWizard | (wave 4 + fix) |

## Wave 5
| BBX-2-02 IIS worker | committed | P | 5 | U: MED out-of-process match + paren parsing fixed | (wave 5 + fix) |
| BBX-3-01 CLR | committed | P | 5 | U: LOW no-target footprint fixed | (wave 5 + fix) |
| BBX-3-02 .NET10 fallback | committed | P | 5 | U: paren parsing + no-target footprint fixed | (wave 5 + fix) |
| BBX-3-03 per-proc I/O | committed | P | 5 | U: paren parsing + no-target footprint fixed | (wave 5 + fix) |

## Wave 6
| BBX-4-01 event log | pending | P | 6 | — | — |
| BBX-4-02 TCP | pending | P | 6 | — | — |
| BBX-4-03 ASP.NET | pending | E | 6 | — | — |
| BBX-4-04 mem pressure | pending | E | 6 | — | — |
| BBX-4-05 busiest shape | pending | P | 6 | — | — |

## Wave 7 (SQL 5a)
| BBX-5-01 SQL service+DACL | pending | P | 7 | — | — |
| BBX-5-03 recorder codes | pending | P | 7 | — | — |
| BBX-5-05 allowlist+lint | pending | P | 7 | — | — |
| BBX-5-08 back-off | pending | P | 7 | — | — |

## Wave 8
| BBX-5-02 SQL CLI+grant | pending | P | 8 | — | — |
| BBX-5-04 /status C11 | pending | P | 8 | — | — |
| BBX-5-06 conn controls | pending | P | 8 | — | — |
| BBX-5-09 output+caps | pending | P | 8 | — | — |

## Wave 9
| BBX-5-07 ceiling+errors | pending | P | 9 | — | — |

## Wave 10 — ULTIMATE GATE (BL-123+BL-124) before samplers
| BBX-5-10 runner+C13 | pending | P | 10 | — | — |

## Wave 11 (samplers)
| BBX-5-11 requests | pending | P | 11 | — | — |
| BBX-5-12 waits | pending | P | 11 | — | — |
| BBX-5-13 schedulers | pending | P | 11 | — | — |
| BBX-5-14 filestats+AG | pending | P | 11 | — | — |

## Wave 12
| BBX-5-16 lane server | pending | P | 12 | — | — |
| BBX-6-02 dump catalogue | pending | P | 12 | — | — |

## Wave 13
| BBX-5-17 lane web | pending | P | 13 | — | — |

## Blocked (written, not dispatched)
| BBX-3-04 EventPipe | blocked | P | — | live pipe proof + cost | — |
| BBX-6-01 VM steal | blocked | P | — | live typeperf on VM | — |
| BBX-6-03 alloc sampling | blocked | P | — | EventPipe proof + cost + owner go | — |

## Log
- 2026-11-13: repo cloned, branch created, docs read, surveys done, Go 1.27.3 installed, baseline green. PLAN.md, 44 packets, GRAPH.md, STATE.md written.
- 2026-11-13: subagents created in .qoder/agents. NOTE: this session's Agent tool only exposes general-purpose/Explore/Plan; custom types not dispatchable and /agents reload not available as a tool. Dispatch uses general-purpose with the tier prompt embedded; per-dispatch model tier cannot be forced from here (agent files carry the intended tier for CLI use).
- 2026-11-13: Wave 1 (1-01..1-04) implemented in worktrees, merged, Ultimate batch review found 1 HIGH (budget slow-group livelock) + 1 LOW; fix round applied (NoteSkipped + GroupTiming optional flag), re-review confirmed resolved. Committed 06add92. CounterReading rename propagated to packets 1-06/1-09; NoteSkipped wiring added to 1-05.
- 2026-11-13: Wave 2 (1-05,1-06,1-10) implemented, merged, integration tests green. Ultimate review found 1 CRITICAL: PdhExpandWildCardPathW called with 4 args (missing dwFlags) → nondeterministic wildcard expansion. Fixed (5th arg 0 at pdh_windows.go:99,107); all PDH call arities re-verified against Win32 signatures; build/vet/tests pass.
- 2026-11-13: Wave 3 (1-07,1-08,1-09,1-12) implemented, merged. Ultimate review found HIGH (CPU self-throttle divided by cores → budget never fires) + MED (widening fabricates running-late/false gaps) + LOW (resolver Start==0 cache) + missing 1-07 report. Fix r1 applied; re-review confirmed but flagged a NEW data race (EffectiveInterval read from HTTP goroutines) + offline-reader false-late. Fix r2: atomic.Int64 effective interval + manifest EffectiveIntervalSeconds used by DeriveState. All blackbox tests pass. -race not runnable here (no gcc); recommend owner runs -race on a CGO host. Server-suite FAIL confirmed pre-existing flake (TestAsyncImportCommitsInTheBackground, fails on pristine tree).
- 2026-11-14: Wave 4 (1-11,2-01,2-03,2-04) implemented, merged. 2-03 blocked by pre-existing clumsy_test race (real-PDH group widens window) → test-only stabilization (processed-tick barrier in hostRig). 2-04 Windows PDH object enumeration wired (sqlhost_windows.go). Integration caught pdhDetailWizard redeclaration (parallel worktrees) → deduped. Ultimate review found HIGH (all-or-nothing pool/service cache) + HIGH (W3WP PDH cache never refreshed → stale on recycle) + MED (sqlhost 2nd PdhEnumObjects bRefresh) + wiring gap (resolver never called in host.go). All fixed (independent caches, refresh-before-items, processChangeGate wiring, bRefresh=0); re-review confirmed; suite ×3 green.
- 2026-11-14: Wave 5 (2-02,3-01,3-02,3-03 per-process groups) implemented, merged; suite ×3 green in isolation. Ultimate review found MED (IIS misses out-of-process .NET10 pools — matched retargeted child PID not w3wp) + LOW-MED (instance parsing drops names with ')') + LOW (clr/proc10/procIO sample machine-wide with no targets). Fixed: IIS matches by pool name (reports w3wp PID), ')\ ' parsing in iis/proc10/procIO, no-target short-circuit. Tests added (out-of-process, paren, no-target-zero-collect). NOTE: pre-existing load-sensitive clumsy listen-timeout flake (TestClumsyWrongReadKeyIsRefusedEveryWay) appears only under parallel test load, passes in isolation.