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

export const BrokerCoordinatorDashboard: Dashboard = {
  variates: [
    {
      tagKey: "node",
      label: "Node",
      db: MonitoringDB,
      multiple: true,
      sql: "show tag values from 'lindb.broker.state_manager' with key=node",
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
                sql: "select 'emit_events' from 'lindb.broker.state_manager' where type='node_joins' group by node",
                watch: ["node"],
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
                sql: "select 'emit_events' from 'lindb.broker.state_manager' where type='node_leaves' group by node",
                watch: ["node"],
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
            title: "Storage State Changed",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'emit_events' from 'lindb.broker.state_manager' where type='storage_changes' group by node",
                watch: ["node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Storage State Deletion",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'emit_events' from 'lindb.broker.state_manager' where type='storage_deletes' group by node",
                watch: ["node"],
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
            title: "Database Config Changed",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'emit_events' from 'lindb.broker.state_manager' where type='database_changes' group by node",
                watch: ["node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Database Config Deletion",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'emit_events' from 'lindb.broker.state_manager' where type='database_deletes' group by node",
                watch: ["node"],
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
            title: "Panic(Process Event)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'panics' from 'lindb.broker.state_manager' group by node",
                watch: ["node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 24,
        },
      ],
    },
  ],
};