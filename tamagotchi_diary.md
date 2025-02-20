# Diary

## Installation

```bash
git clone https://github.com/Argus-Labs/world-engine.git && cd world-engine && git checkout v8.0.2
go install github.com/golang/mock/mockgen@v1.6.0
make generate
make rollup-proto-gen
make rollup-build
```

I have to remove target: nakama in docker-compose.yml

### Get celestia devnet token

```shell
docker compose -f ./docker-compose.yml build
make start-da
```

###

```shell
export DA_AUTH_TOKEN=$(docker exec $(docker ps -q) /bin/sh -c "celestia bridge auth admin --node.store /home/celestia/bridge")
export DA_AUTH_TOKEN="..."
docker compose up
```

### Extra

```bash
export DA_AUTH_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdfQ.yAIW8zw75eYOdbTbI6MDk9bULNsJ_5e2h4KKUGq2VXY"
make start-sequencer
```

```bash
export DA_AUTH_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdfQ.yAIW8zw75eYOdbTbI6MDk9bULNsJ_5e2h4KKUGq2VXY"
make game
```

### javascript

[Account](https://github.com/Argus-Labs/cardinal-editor/blob/f125afaf7cac8be09d0cb437b564cf24c86902df/src/lib/account.ts)

```shell
curl --request POST \
  --url http://localhost:4040/tx/game/create-player \
  --header 'Content-Type: application/json' \
  --data '{
    "personaTag": "CoolMage",
    "namespace": "",
    "nonce": 0,
    "signature": "",
    "body": {
    "nickname": "CoolMage"
    }
  }'
```

## Flowchart

### Architecture

```mermaid
---
title: Frontend / Abstract Accounts / Cardinal
---
graph TD

    User["User"]
    Email["Email"]

    subgraph Frontend["web3"]
        subgraph WC["World Client"]
            WC_A["nakama-js"]
            subgraph Session["User Session"]
                authToken
                Persona
            end
            WC_B["Rest API"] 
        end
        Proxy["Web Proxy/Balancer"]
        UI["Web Interface"]

        Proxy --> UI
        UI --> WC
        WC_A --> Session
    end

    subgraph nakamaCluster["Nakama Cluster"]
        subgraph nakama[Node]
        end
        Account["Account abstraction"]
        cockroachdb[(CockroachDB)]

        nakama --> cockroachdb
    end

    subgraph cardinal[Cardinal Service]
    end
    
    WC_B <-- HTTP GET/POST --> cardinal
    cardinal <-- CARDINAL_NAMESPACE --> nakama

    %% Abstract Account
    User -- Login --> Proxy
    WC_A -- Authentication --> nakamaCluster
    nakamaCluster --> Session
    nakama --> Account
    Account <--> Email

    classDef green fill:#696,stroke:#333;
    classDef database fill:#f0f0f0;
    classDef service fill:#f0fff0;

    class cockroachdb database;
    class cardinal green;
    class nakamaCluster service;
```

```mermaid
---
title: Cardinal / EVM / DAS
---
graph TD

    subgraph celestia["Celestia Data Availability"]
        CL_LN["Celestia Light Node"]
        CL_BN["Celestia Bridge Node"] 
        CL_FN["Celestia Full Storage Node"]

        CL_LN-- "Store Data Call" --> CL_BN
        CL_BN-- "celestia-app" -->CL_FN
    end

    subgraph chain[EVM Chain]
        subgraph sequencer
            grpc["Grpc:9061"]
        end
        subgraph sc[smart contracts]
            faucet[Faucet]
            router[Router]
            erc20[Native Token]
            payment["Payment logic"]

            payment <--> erc20
            faucet -- mint --> erc20
        end
    end

    subgraph cardinal[cardinal]
        subgraph SCI["Shard Comm Int"]
            ShardRouter["Router"]
        end
        subgraph State["State"]
            ECS["ECS"]
        end
    end
    
    redis[(Cache)] <--> cardinal
    celestia <-- DA_AUTH_TOKEN, DA_NAMESPACE_ID --> chain
    router <-- BASE_SHARD_ROUTER_KEY --> ShardRouter

    Frontend <-- 1- HTTP GET/POST --> cardinal
    ShardRouter -- 2- sendMessage[PersonaTag,Data] --> router
    router -- 3- Pay --> payment
    payment -- 4- Result --> router
    router -- 5- messageResult --> ShardRouter

    ECS -- Txs --> grpc
    grpc -- Txs --> celestia

    %% Network Connections
    classDef database fill:#f0f0f0;
    classDef blockchain fill:#e8f5e9;
    classDef service fill:#f0fff0;
    classDef monitoring fill:#f0d9ff;

    class cockroachdb database;
    class redis database;
    class celestia-devnet blockchain;
    class chain blockchain;
    class game service;
    class nakama service;
    class jaeger monitoring;
    class test_nakama service;
```

## Explanation of the Diagram

### Core Services

- **CockroachDB**: Acts as the primary database for both the `game` and `nakama` services.
- **Redis**: Provides in-memory data storage for the `game` and `nakama` services.

### Blockchain Services

- **Celestia Devnet**: Serves as the blockchain node for the `chain` service.
- **EVM Chain**: Interacts with the `game` service to enable blockchain functionality.

### Game Services

- **Game Service**: Connects to `nakama`, `redis`, and `jaeger` for gameplay, data storage, and monitoring.
- **Nakama**: A server for handling game logic, connected to `cockroachdb`, `redis`, and `jaeger`.

### Monitoring & Tracing

- **Jaeger**: Collects tracing data from the `game`, `nakama`, and `chain` services for performance monitoring.

### Test Services

- **Test Nakama**: A test service that depends on the `nakama` server.

## Design

## Entity

- Pet
- LeaderBoard
- Food
  - Berrie
  - Snack
  - Yum
  - Pop
  - Honey
  - Munch
  - leaf
- Toys
  - Ball
  - Stick
  - Jump rope
- Habitats
  - park (walk)
  - forest (run)
  - beach (swim)
- care
  - injection
  - vitamins

## Component

- energia (dormir, comer sube, jugar baja) baja con el tiempo
- vida (enferma: baja, hambre 100%: baja, comer: sube, curar: deja de bajar, pero no sube). Si duerme
- wellness: sube jugando y es la propiedad para el ranking.
- higiene:
- DNA:

estado: 
    - baniarse: reduce posibilidad de enfermedad y sube higiene. Jugar baja higiene.
    - enfermo -> baja vida. Con injeccion cura, detiene el danio constante de enfermedad.
    - vitaminas: sube vida rapidamente.
    - comida: sube energia.
    - energia: se usa para juegar, el cual sube el wellness
    - dormir: sube energia, sube vida lentamente.
    - juguetes: sube rapido el wellness.