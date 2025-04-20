<script>
import {BsPersonArmsUp, FaDollarSign} from "@kalimahapps/vue-icons";
import { AkCircleTriangleLeft } from '@kalimahapps/vue-icons';
export default {
  name: "PaymentsPage",
  components: {BsPersonArmsUp, FaDollarSign,AkCircleTriangleLeft},
  methods: {
    goToMain() {
      this.$router.push('/')
    }
  }
}
</script>

<template>
  <div class="container">
    <div class="row" style="justify-content: left">
      <div class="col-1">
        <button class="btn" v-on:click="goToMain">Назад</button>
      </div>
    </div>
    <div class="row">
      <div class="col-1">
          <h2>Профиль</h2>
          <p>Информация о вашем аккаунте</p>
      </div>
    </div>
  </div>
</template>