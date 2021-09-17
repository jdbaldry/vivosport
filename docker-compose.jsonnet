{
  version: '3.1',

  local db = 'postgres',
  local user = 'vivosport',
  local password = 'vivosport',
  services: {
    [db]: {
      image: db,
      restart: 'always',

      environment: {
        POSTGRES_DB: user,
        POSTGRES_PASSWORD: password,
        POSTGRES_USER: user,
      },
      ports: ['${POSTGRES_PORT:-5432}:5432'],
      volumes: ['/home/jdb/pgdata:/var/lib/postgresql/data:z'],
    },
    grafana: {
      image: 'grafana/grafana:7.3.3',

      entrypoint:
        local dashboards = {
          apiVersion: 1,
          providers: [{
            name: 'vivosport',

            allowUiUpdates: true,
            options: { path: '/var/lib/grafana/dashboards' },
            updateIntervalSeconds: 1,
          }],
        };
        local dataSources = {
          apiVersion: 1,
          datasources: [{
            name: 'vivosport',

            database: user,
            url: db,  // host
            jsonData: {
              sslmode: 'disable',
            },
            isDefault: true,
            secureJsonData: {
              password: password,
            },
            type: db,
            user: user,
          }],
        };
        [
          'sh',
          '-euc',
          |||
            printf "%s" > /etc/grafana/provisioning/dashboards/vivosport.yml
            printf "%s" > /etc/grafana/provisioning/datasources/vivosport.yml
            exec /run.sh
          ||| % std.map(std.manifestYamlDoc, [dashboards, dataSources]),
        ],
      depends_on: [db],
      environment: [
        'GF_AUTH_ANONYMOUS_ENABLED=true',
        'GF_AUTH_ANONYMOUS_ORG_ROLE=Admin',
        'GF_DATABASE_TYPE=' + db,
        'GF_DATABASE_HOST=' + db,
        'GF_DATABASE_USER=' + user,
        'GF_DATABASE_PASSWORD=' + password,
        'GF_DATABASE_SSL_MODE=disable',
      ],
      ports: ['${GRAFANA_PORT:-3000}:3000'],
    },
  },
}
