import requests

BASE_URL = "http://localhost:8080/api"

def test_send_coins():
    auth_data_u1 = {
        "username": "user001", 
        "password": "password"
    }
    auth_response = requests.post(f"{BASE_URL}/auth", json=auth_data_u1)
    assert auth_response.status_code == 200
    token1 = auth_response.json().get("token")
    headers1 = {"Authorization": f"Bearer {token1}"}

    user_before_response = requests.get(f"{BASE_URL}/info", headers=headers1)
    assert user_before_response.status_code == 200
    balance_before = user_before_response.json().get("coins")

    auth_data_u2 = {
        "username": "user002", 
        "password": "password"
    }
    auth_response = requests.post(f"{BASE_URL}/auth", json=auth_data_u2)
    assert auth_response.status_code == 200
    token2 = auth_response.json().get("token")
    headers2 = {"Authorization": f"Bearer {token2}"}
    amount = 100
    send_data = {"toUser": "user002", "amount": amount}
    send_response = requests.post(f"{BASE_URL}/sendCoin", json=send_data, headers=headers1)
    assert send_response.status_code == 200

    user_after_response = requests.get(f"{BASE_URL}/info", headers=headers1)
    assert user_after_response.status_code == 200
    balance_after = user_after_response.json().get("coins")

    assert amount == balance_before - balance_after

def test_buy_merch():
    auth_data = {
        "username": "user003", 
        "password": "password"
    }
    auth_response = requests.post(f"{BASE_URL}/auth", json=auth_data)
    assert auth_response.status_code == 200
    token = auth_response.json().get("token")
    headers = {"Authorization": f"Bearer {token}"}

    user_before_response = requests.get(f"{BASE_URL}/info", headers=headers)
    assert user_before_response.status_code == 200
    balance_before = user_before_response.json().get("coins")

    item = "cup"
    buy_response = requests.get(f"{BASE_URL}/buy/{item}", headers=headers)
    assert buy_response.status_code == 200

    user_after_response = requests.get(f"{BASE_URL}/info", headers=headers)
    assert user_after_response.status_code == 200
    balance_after = user_after_response.json().get("coins")

    assert balance_after < balance_before