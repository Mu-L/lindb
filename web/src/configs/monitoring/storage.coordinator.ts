/*
Licensed to LinDB under one or more contributor
license agreements. See the NOTICE file distributed with
this work for additional information regarding copyright
ownership. LinDB licenses this file to you under
the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
 
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
import { MonitoringDB } from "@src/constants";
import { Dashboard, UnitEnum } from "@src/models";

export const StorageCoordinatorDashboard: Dashboard = {
  variates: [
    {
      tagKey: "namespace",
      label: "Namespace",
      db: MonitoringDB,
      sql: "show tag values from 'lindb.storage.state_manager' with key=namespace",
      watch: { clear: ["node"] },
    },
    {
      tagKey: "node",
      label: "Node",
      watch: { cascade: ["namespace"] },
      db: MonitoringDB,
      multiple: true,
      sql: "show tag values from 'lindb.storage.state_manager' with key=node",
    },
  ],
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Node Joins",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'emit_events' from 'lindb.storage.state_manager' where type='node_joins' group by namespace,node",
                watch: ["namespace", "node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Node Leaves",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'emit_events' from 'lindb.storage.state_manager' where type='node_leaves' group by namespace,node",
                watch: ["namespace", "node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Shard Assigns",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'emit_events' from 'lindb.storage.state_manager' where type='shard_assigns' group by namespace,node",
                watch: ["namespace", "node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Panic(Process Event)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'panics' from 'lindb.storage.state_manager' group by namespace,node",
                watch: ["namespace", "node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
      ],
    },
  ],
};
