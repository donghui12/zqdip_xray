import json
import os
from flask import Flask, request, jsonify

app = Flask(__name__)

# Xray 配置文件路径
CONFIG_PATH = "/usr/local/etc/xray/config.json"


# 读取配置文件
def read_config():
    if os.path.exists(CONFIG_PATH):
        with open(CONFIG_PATH, 'r') as f:
            return json.load(f)
    return {}


# 写入配置文件
def write_config(config):
    with open(CONFIG_PATH, 'w') as f:
        json.dump(config, f, indent=4)


# 新增路由规则
@app.route('/v1/server/add', methods=['POST'])
def add_rule():
    data = request.get_json()
    inbound_tag = data.get('inboundTag')
    outbound_tag = data.get('outboundTag')

    if not inbound_tag or not outbound_tag:
        return jsonify({"error": "inboundTag and outboundTag are required"}), 400

    config = read_config()
    rules = config.get('routing', {}).get('rules', [])

    # 检查是否已存在相同的 inboundTag
    for rule in rules:
        if inbound_tag in rule.get('inboundTag', []):
            return jsonify({"error": "InboundTag already exists"}), 400

    new_rule = {
        "type": "field",
        "inboundTag": [inbound_tag],
        "outboundTag": outbound_tag
    }

    rules.append(new_rule)
    config['routing']['rules'] = rules
    write_config(config)

    return jsonify({"inboundTag": inbound_tag, "outboundTag": outbound_tag}), 201


# 删除规则、入站、出站同时执行
@app.route('/v1/server/remove', methods=['POST'])
def remove_rule_inbounds_outbounds():
    data = request.get_json()
    inbound_tag = data.get('inboundTag')

    if not inbound_tag:
        return jsonify({"error": "inboundTag is required"}), 400

    config = read_config()

    # 删除 routing.rules 中的规则
    rules = config.get('routing', {}).get('rules', [])
    rules = [rule for rule in rules if inbound_tag not in rule.get('inboundTag', [])]
    config['routing']['rules'] = rules

    # 删除 inbounds 中的规则
    inbounds = config.get('inbounds', [])
    inbounds = [inbound for inbound in inbounds if inbound_tag != inbound.get('tag')]
    config['inbounds'] = inbounds

    # 删除 outbounds 中的规则
    outbounds = config.get('outbounds', [])
    outbounds = [outbound for outbound in outbounds if inbound_tag != outbound.get('tag')]
    config['outbounds'] = outbounds

    write_config(config)

    return jsonify({"message": "Rule, Inbound, and Outbound removed"}), 200


# 更新 inbounds 中的 port、用户名、密码
@app.route('/v1/server/update', methods=['POST'])
def update_inbounds():
    data = request.get_json()
    inbound_tag = data.get('inboundTag')
    new_port = data.get('port')
    new_user = data.get('user')
    new_pass = data.get('pass')

    if not inbound_tag or not new_port or not new_user or not new_pass:
        return jsonify({"error": "inboundTag, port, user, and pass are required"}), 400

    config = read_config()
    inbounds = config.get('inbounds', [])

    for inbound in inbounds:
        if inbound.get('tag') == inbound_tag:
            inbound['port'] = new_port
            inbound['settings']['accounts'][0]['user'] = new_user
            inbound['settings']['accounts'][0]['pass'] = new_pass
            write_config(config)
            return jsonify({"message": "Inbound updated"}), 200

    return jsonify({"error": "InboundTag not found"}), 404



# 获取路由规则列表
@app.route('/v1/server/list', methods=['GET'])
def list_rules():
    config = read_config()
    rules = config.get('routing', {}).get('rules', [])
    print(rules)
    rule_list = []
    for rule in rules:
        if rule.get("inboundTag"):
            rule_list.append(rule['inboundTag'][0])
    return jsonify({"节点": rule_list}), 200


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8888)