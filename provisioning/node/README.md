# Node

## Discovery / Alive monitoring

- With [memberlist](https://github.com/hashicorp/memberlist).

### Agentが正常開始 (Node作成 -> memberlist参加)

1. AgentがAPIにNodeを保存、APIはmemberlistにないのでNotReadyに
2. AgentがAPIを起点にmemberlistに参加、APIは通知からNodeをReadyに

### Agentが正常開始 (memberlist参加 -> Node作成)

1. AgentがAPIを起点にmemberlistに参加、APIはNodeがないのでなにも変更せず
2. AgentがAPIにNodeを保存、APIのmemberlistにあるのでReadyに

### Agentが正常終了

1. AgentがAPIからNodeを削除
2. Agentがmemberlistから抜ける

### Agentが異常終了

1. Agentがmemberlistから抜ける
2. APIが離脱を検知し、NodeをNotReadyに

### APIが異常終了 (同時にAgentが異常終了した場合も同様)

- APIが死亡時は動作としては問題がない

#### APIが復活

- TODO: AgentとAPIどちらがどちらのmemberlistにジョインするか
  - Agentからな気がする
