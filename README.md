# BoardGameTracker API

A simple REST API for keeping track of boardgame stats.

## Endpoints

##### GET `/boardGames`
Searching boardgames

**Parameters**

|          Name | Required |  Type   | Description                                                                                                                                                           |
| -------------:|:--------:|:-------:| -------------- |
|     `search` | required | string  | The search term    |
|     `limit` | optional | int  | Maximum results. Default: 100 |
|     `room` | optional | int  | The id of the room to search for prioritizing played games |

**Response**
```json
[  
   {  
      "id": 16933,
      "name":"Monopoly (1933)"
   }
]
```
---
##### GET `/room/{room}`
Getting room information

**Parameters**

|          Name | Required |  Type   | Description                                                                                                                                                           |
| -------------:|:--------:|:-------:| -------------- |
|     `room` | required | string  | The room name    |

**Response**
```
{  
   "id":1,
   "hash":"098f6bcd4621d373cade4e832627b4f6",
   "createTime":"2019-03-09T11:25:18Z",
   "players":[  
      {  
         "id":1,
         "joinTime":"2019-03-10T13:11:48Z",
         "name":"Tom",
         "color":"#ffffff"
      },
      {  
         "id":2,
         "joinTime":"2019-03-14T16:20:26Z",
         "name":"Jerry",
         "color":"#ff0000"
      }
   ],
   "games":[  
      {  
         "id":1,
         "boardGame":{  
            "id":98238,
            "name":"Mousetrap Trappers (2015)"
         },
         "gameTime":"2019-03-31T14:38:40Z",
         "gamePlayers":[  
            {  
               "player":{  
                  "id":1,
                  "joinTime":"2019-03-14T16:20:35Z",
                  "name":"Tom",
                  "color":"#ff0000"
               },
               "points":100
            },
            {  
               "player":{  
                  "id":2,
                  "joinTime":"2019-03-14T16:20:40Z",
                  "name":"Jerry",
                  "color":"#ff0000"
               },
               "points":1
            }
         ]
      }
   ]
}
```

---
##### POST `/room/{room}/player/`

**Parameters**

|          Name | Required |  Type   | Description                                                                                                                                                           |
| -------------:|:--------:|:-------:| -------------- |
|     `room` | required | string  | The room name    |

**Body**
```
{  
    "name": string,
    "color": string
}
```

**Response**

```json
{  
    "id":1,
    "joinTime":"2019-03-14T16:20:35Z",
    "name":"Tom",
    "color":"#ff0000"
}
```

---

##### POST `/room/{room}/game/`

**Parameters**

|          Name | Required |  Type   | Description                                                                                                                                                           |
| -------------:|:--------:|:-------:| -------------- |
|     `room` | required | string  | The room name    |

**Body**
```
{  
    "boardGame": {
      "id": number
    },
    "gamePlayers": [
        {
            "points": number,
            "player": {
                "id": number
            }
        }
    ]
}
```

**Response**

```json
{
 "id":1,
 "boardGame":{  
    "id":98238,
    "name":"Mousetrap Trappers (2015)"
 },
 "gameTime":"2019-03-31T14:38:40Z",
 "gamePlayers":[  
    {  
       "player":{  
          "id":1,
          "joinTime":"2019-03-14T16:20:35Z",
          "name":"Tom",
          "color":"#ff0000"
       },
       "points":100
    },
    {  
       "player":{  
          "id":2,
          "joinTime":"2019-03-14T16:20:40Z",
          "name":"Jerry",
          "color":"#ff0000"
       },
       "points":1
    }
 ]
}
```


## Database structure

### BoardGame
```
+--------------+---------------------+------+-----+---------+----------------+
| Field        | Type                | Null | Key | Default | Extra          |
+--------------+---------------------+------+-----+---------+----------------+
| boardgame_id | bigint(20) unsigned | NO   | PRI | NULL    | auto_increment |
| name         | varchar(255)        | YES  | UNI | NULL    |                |
+--------------+---------------------+------+-----+---------+----------------+
```

### Game
```
+--------------+---------------------+------+-----+-------------------+----------------+
| Field        | Type                | Null | Key | Default           | Extra          |
+--------------+---------------------+------+-----+-------------------+----------------+
| game_id      | bigint(20) unsigned | NO   | PRI | NULL              | auto_increment |
| room_id      | bigint(20) unsigned | NO   | MUL | NULL              |                |
| boardgame_id | bigint(20) unsigned | NO   | MUL | NULL              |                |
| game_time    | timestamp           | NO   |     | CURRENT_TIMESTAMP |                |
+--------------+---------------------+------+-----+-------------------+----------------+
```

### GamePlayers
```
+---------------+---------------------+------+-----+---------+----------------+
| Field         | Type                | Null | Key | Default | Extra          |
+---------------+---------------------+------+-----+---------+----------------+
| gameplayer_id | bigint(20) unsigned | NO   | PRI | NULL    | auto_increment |
| game_id       | bigint(20) unsigned | NO   | MUL | NULL    |                |
| player_id     | bigint(20) unsigned | NO   | MUL | NULL    |                |
| points        | int(11)             | YES  |     | NULL    |                |
+---------------+---------------------+------+-----+---------+----------------+
```
### Player
```
+-----------+---------------------+------+-----+-------------------+----------------+
| Field     | Type                | Null | Key | Default           | Extra          |
+-----------+---------------------+------+-----+-------------------+----------------+
| player_id | bigint(20) unsigned | NO   | PRI | NULL              | auto_increment |
| room_id   | bigint(20) unsigned | NO   | MUL | NULL              |                |
| name      | varchar(255)        | NO   |     | <empty>           |                |
| color     | varchar(7)          | NO   |     | #7dbc67           |                |
| join_time | timestamp           | NO   |     | CURRENT_TIMESTAMP |                |
+-----------+---------------------+------+-----+-------------------+----------------+
```

### Room

```
+-------------+---------------------+------+-----+-------------------+----------------+
| Field       | Type                | Null | Key | Default           | Extra          |
+-------------+---------------------+------+-----+-------------------+----------------+
| room_id     | bigint(20) unsigned | NO   | PRI | NULL              | auto_increment |
| room_hash   | char(32)            | NO   | MUL | NULL              |                |
| create_time | timestamp           | NO   |     | CURRENT_TIMESTAMP |                |
+-------------+---------------------+------+-----+-------------------+----------------+
```
