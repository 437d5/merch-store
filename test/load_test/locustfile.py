from locust import HttpUser, task, between
import random
import string

item_names = ["t-shirt", "cup", "book", "pen", "powerbank", "hoody", "umbrella",
                "socks", "wallet", "pink-hoody"]

def generate_username(index):
    return f"user{index:03d}"

password = ["password" for i in range(1, 100001)]

users = [{"username":generate_username(i), "password":password[i-1]} for i in range(1, 100001)]

class UserBehavior(HttpUser):
    wait_time = between(1, 2)  # Задержка между запросами 1-2 сек
    idx = 0
    def on_start(self):
        self.user_data = users[UserBehavior.idx]
        UserBehavior.idx += 1
        response = self.client.post("/api/auth", json=self.user_data)

        if response.status_code == 200:
            self.token = response.json().get("token")
        else:
            self.token = None
            print(f"Ошибка входа для {self.user_data['username']}")

    @task(3)
    def get_user_info(self):
        """Тест запроса информации о пользователе"""
        if self.token:
            headers = {"Authorization": f"Bearer {self.token}"}
            self.client.get("/api/info", headers=headers)

    @task(2)
    def buy_item(self):
        """Тест покупки предмета"""
        if self.token:
            headers = {"Authorization": f"Bearer {self.token}"}
            item_name = random.choice(item_names)
            self.client.get(f"/api/buy/{item_name}", headers=headers)

    @task(1)
    def send_coins(self):
        """Перевод монет другому пользователю"""
        if self.token:
            headers = {"Authorization": f"Bearer {self.token}"}
            recipient = random.choice(users)["username"]
            self.client.post("/api/sendCoin", json={"toUser": recipient, "amount": 100}, headers=headers)