<template>
  <div>
    <p v-if="!isDefault">
      <label for="username">{{ $t("settings.username") }}</label>
      <input
        class="input input--block"
        type="text"
        v-model="user.username"
        id="username"
        :disabled="!this.isNew"
      />
    </p>

    <p v-if="!isDefault">
      <label for="password">{{ $t("settings.password") }}</label>
      <input
        :class="passwordClass"
        type="password"
        :placeholder="passwordPlaceholder"
        v-model="user.password"
        id="password"
      />
      <input
          :class="passwordClass"
          type="password"
          :placeholder="$t('settings.newPasswordConfirm')"
          v-model="passwordConf"
          name="password"
      />
  </p>

    <p>
      <label for="scope">{{ $t("settings.scope") }}</label>
      <input
        class="input input--block"
        type="text"
        v-model="user.scope"
        id="scope"
      />
    </p>

    <p>
      <label for="locale">{{ $t("settings.language") }}</label>
      <languages
        class="input input--block"
        id="locale"
        :locale.sync="user.locale"
      ></languages>
    </p>

    <p v-if="!isDefault">
      <input
        type="checkbox"
        :disabled="user.perm.admin"
        v-model="user.lockPassword"
      />
      {{ $t("settings.lockPassword") }}
    </p>

    <permissions :perm.sync="user.perm" />
    <commands v-if="isExecEnabled && isCmdLimitEnabled" :commands.sync="user.commands" />

    <div v-if="!isDefault && false">
      <h3>{{ $t("settings.rules") }}</h3>
      <p class="small">{{ $t("settings.rulesHelp") }}</p>
      <rules :rules.sync="user.rules" />
    </div>
  </div>
</template>

<script>
import Languages from "./Languages";
import Rules from "./Rules";
import Permissions from "./Permissions";
import Commands from "./Commands";
import { enableExec } from "@/utils/constants";
import { enableCmdLimit } from "@/utils/constants";

export default {
  name: "user",
  data: function () {
    return {      
      passwordConf: "",
    };
  },
  components: {
    Permissions,
    Languages,
    Rules,
    Commands,
  },
  props: ["user", "isNew", "isDefault"],
  computed: {
    passwordPlaceholder() {
      return this.isNew ? "" : this.$t("settings.avoidChanges");
    },
    isExecEnabled: () => enableExec,
    isCmdLimitEnabled: () => enableCmdLimit,
    passwordClass() {
      const baseClass = "input input--block";

      if ((this.user.password === "" || this.user.password === undefined) && this.passwordConf === "") {
        return baseClass;
      }
             
      var reg = /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$%^&*-]).{8,}$/;
      var hangulcheck = /[ㄱ-ㅎ|ㅏ-ㅣ|가-힣]/;
      
      if(false === reg.test(this.user.password)) {
        return `${baseClass} input--red`;
      }else if(/(\w)\1\1\1/.test(this.user.password)){
       return `${baseClass} input--red`;
      }else if(this.user.password.search(/\s/) != -1){
        return `${baseClass} input--red`;
      }else if(hangulcheck.test(this.user.password)){
       return `${baseClass} input--red`;
      }else {
        if (this.user.password !== this.passwordConf || this.user.password === "") {
          return `${baseClass} input--red`;
        }
      }

      if (this.user.password === this.passwordConf) {
        return `${baseClass} input--green`;
      }

      return `${baseClass} input--red`;
    },
  },
  watch: {
    "user.perm.admin": function () {
      if (!this.user.perm.admin) return;
      this.user.lockPassword = false;
    },
  },
};
</script>
