<!--
SPDX-FileCopyrightText: 2022-present Intel Corporation

SPDX-License-Identifier: Apache-2.0
-->


# Metric Function

![Metric Function Architecture](/docs/images/Metric_Function_Arch.png)



# Supported Features
1. API Server exposure
2. Prometheus Client exposure
3. Analytics Function exposure(not supported in this release)

# Types of Statistics
1. Core Subscriber information
2. Network Function Status(only UPF and GNodeB supported)
3. Core Message Statistics(only SMF and AMF supported)

# API Server APIs supported
1. GetSubscriberSummary (/nmetric-func/v1/subscriber/<imsi>)
2. GetSubscriberAll (/nmetric-func/v1/subscriber/all)
3. GetNfStatus (/nmetric-func/v1/nfstatus/<GNB/UPF>)
4. GetNfServiceStats (/nmetric-func/v1/nfServiceStatsSummary/<AMF/SMF>)
5. GetNfServiceStatsAll (/nmetric-func/v1/nfServiceStats/all)

For more details about the Grafana Dashboard, please refer- https://docs.aetherproject.org/master/developer/aiabhw5g.html#enable-monitoring

  
