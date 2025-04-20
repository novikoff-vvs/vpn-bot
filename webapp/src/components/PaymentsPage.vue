<script>
import {BsPersonArmsUp, FaDollarSign, McLoading2Line} from "@kalimahapps/vue-icons";

export default {
  name: "PaymentsPage",
  components: {BsPersonArmsUp, FaDollarSign, McLoading2Line},
  data: () => {
    {
      return {
        is_loading: false
      }
    }
  },
  methods: {
    goToMain() {
      this.$router.push('/')
    },
    getTelegramChatId() {
      // Проверяем, доступен ли объект Telegram WebApp
      if (window.Telegram && window.Telegram.WebApp) {
        try {
          // Пытаемся получить данные initData
          const initData = window.Telegram.WebApp.initData || '';

          // Если есть initData, пробуем извлечь chat_id
          if (initData) {
            const params = new URLSearchParams(initData);
            const chatId = params.get('chat_id');

            if (chatId) {
              return chatId;
            } else {
              console.warn('chat_id не найден в initData');
              return null;
            }
          } else {
            console.warn('initData отсутствует в Telegram WebApp');
            return null;
          }
        } catch (error) {
          console.error('Ошибка при получении chat_id:', error);
          return null;
        }
      } else {
        console.warn('Telegram WebApp не доступен');
        return null;
      }
    },
    async sendTelegramData() {
      try {
        this.is_loading = true
        // Получаем chat_id
        const chatId = this.getTelegramChatId(); // Используем метод из предыдущего примера

        if (!chatId) {
          console.error('Не удалось получить chat_id');
          throw new Error('Не удалось получить chat_id');
        }

        // Данные для отправки
        const postData = {
          chat_id: chatId,
          user_data: window.Telegram.WebApp.initDataUnsafe.user || {},
          // Добавьте другие необходимые данные
        };


        const response = await fetch('vpn-bot:8080/api/webhook/handle', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            // Другие необходимые заголовки
          },
          body: JSON.stringify(postData)
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        console.log('Успешный ответ:', data);
        this.is_loading = false
        return data;

      } catch (error) {
        this.is_loading = false
        console.error('Ошибка при отправке данных:', error);
      } finally {
        this.is_loading = false
      }
    },
  }
}
</script>

<template>
  <div class="container">
    <div class="row">
      <div class="col-1">
        <button class="btn" v-on:click="goToMain">Назад</button>
      </div>
    </div>
    <div class="row">
      <div class="col-1">
        <h2>ОПЛАТА</h2>
      </div>
    </div>
    <div class="row">
      <div class="col-1">
        <div class="card">
          <div class="row">
            <div class="col-1">
              <p>Здесь вы можете произвести оплату вашей подписки.</p>
            </div>
          </div>
          <br>
          <div class="row">
            <div class="col-1">
              <button class="btn" v-on:click="sendTelegramData"
                      style="width: 200px; min-height: 100px; font-size: 20px">
                <p v-if="!is_loading">
                  Оплатить
                </p>
                <McLoading2Line v-else class="animate-spin" style="font-size: 40px"/>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>


<style>
.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>