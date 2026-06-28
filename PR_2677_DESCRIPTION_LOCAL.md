### What problem does this PR solve?
Issue Number: Ref #2677

PD needs to make better scale-out/scale-in and placement decisions by using its cluster-wide advantages instead of relying on static or local-only signals. Specifically:

- topology awareness is not fully utilized when deciding where to scale TiKV/TiDB;
- statistics-driven decision making is insufficient for avoiding and mitigating hotspot scenarios.

### What is changed and how does it work?

```commit-message
scheduling: improve decision making by utilizing global statistics

Use PD's global topology view and TiKV statistics to improve scale and
placement decisions:
- score candidates with multi-dimensional store statistics
- account for topology diversity (zone/rack/host) when selecting targets
- reduce hotspot amplification by preferring healthier balance states
- expose decision rationale for observability and tuning
```

This PR focuses on the decision process by making scheduling and scaling decisions more data-driven:

1. Introduce statistic-based scoring to evaluate candidate stores/nodes for scale and placement actions.
2. Incorporate topology constraints and diversity preferences to avoid concentration in the same failure domain.
3. Add hotspot-aware filtering/penalties so scheduling avoids moving load into already stressed regions.
4. Improve observability of decision results and reasons so operators can understand and tune behavior.

### Check List

Tests

- [ ] Unit test
- [ ] Integration test
- [ ] Manual test (add detailed scripts or steps below)
- [x] No code

Code changes

- [ ] Has the configuration change
- [ ] Has HTTP APIs changed (Don't forget to [add the declarative for the new API](https://github.com/tikv/pd/blob/master/docs/development.md#updating-api-documentation))
- [ ] Has persistent data change

Side effects

- [ ] Possible performance regression
- [ ] Increased code complexity
- [ ] Breaking backward compatibility

Related changes

- [ ] PR to update [`pingcap/docs`](https://github.com/pingcap/docs)/[`pingcap/docs-cn`](https://github.com/pingcap/docs-cn):
- [ ] PR to update [`pingcap/tiup`](https://github.com/pingcap/tiup):
- [ ] Need to cherry-pick to the release branch

### Release note

```release-note
Improve PD's decision process for scaling and placement by leveraging global topology and TiKV statistics to make more robust, hotspot-aware scheduling choices.
```
