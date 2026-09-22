<!--
SPDX-FileCopyrightText: 2022-present Intel Corporation

SPDX-License-Identifier: Apache-2.0
-->
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/omec-project/metricfunc/badge)](https://scorecard.dev/viewer/?uri=github.com/omec-project/metricfunc)

# Metric Function

![Metric Function Architecture](/docs/images/Metric_Function_Arch.png)



# Supported Features
1. API Server exposure
2. Prometheus Client exposure
3. Analytics Function exposure(not supported in this release)

# Types of Statistics
1. Core Subscriber information
2. Network Function Status(only UPF and GNodeB supported)

# API Server APIs supported
1. GetSubscriberSummary (/nmetric-func/v1/subscriber/<imsi>)
2. GetSubscriberAll (/nmetric-func/v1/subscriber/all)
3. GetNfStatus (/nmetric-func/v1/nfstatus/<GNB/UPF>)
4. GetNfStatusAll (/nmetric-func/v1/nfstatus/all)


For more details about the Grafana Dashboard, please refer- https://docs.aetherproject.org/developer/monitoring.html

For more details about the Metric-Function, please refer- https://docs.sd-core.opennetworking.org/master/design/design-metricfunc.html

# Reach out to us through

1. #sdcore-dev channel in [Aether Community Slack](https://aether5g-project.slack.com)
