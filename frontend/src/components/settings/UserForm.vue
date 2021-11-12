<template>
  <div>
    <p>
       <b v-if="user.perm.admin"> [ {{ $t("settings.tgmaccount") }} ]</b>
       <b v-else> [ {{ $t("settings.linuxaccount") }} ]</b>
    </p>
    <p v-if="!isDefault">
      <label for="username">{{ $t("settings.username") }}</label>
      <input
        :class="userNameClass"
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
        name="password"
      />
      <input
          :class="passwordSubClass"
          type="password"
          :placeholder="$t('settings.newPasswordConfirm')"
          v-model="passwordConf"
          name="passwordConf"
      />
  </p>
    
<p v-if="user.perm.admin == false">
   <label for="shell">Shell - 작업중</label>
 </p>

<p v-if="user.perm.admin == false">
  <label for="group">Group - 작업중</label>
 </p>

 <p v-if="user.perm.admin == false">
  <label>계정 유효 일자 - 작업중</label>
 </p>

  <p v-if="user.perm.admin == false">
  <label>암호 기간 만료 설정 - 작업중</label>
 </p>

<p>
      <label for="scope">{{ $t("settings.scope") }}</label>
      <input
        :class="scopeClass"
        type="text"
        v-model="user.scope"
        id="scope"
        :disabled="!this.isNew && $route.path != '/settings/global'"
        style="display:inline-block;"
      /><span v-if="this.isNew"><span v-if="user.username">/</span>{{user.username}}</span>
    </p>

  <p v-if="user.perm.admin == false">
  <label>Quota - 작업중</label>
 </p>

    <p>
       <b v-if="user.perm.admin">&nbsp;</b>
       <b v-else> [ {{ $t("settings.tgmaccount") }} ]</b>
    </p>
    <p>
      <label for="locale">{{ $t("settings.language") }}</label>
      <languages
        class="input input--block"
        id="locale"
        :locale.sync="user.locale"
      ></languages>
    </p>
    <p v-if="!isDefault" v-show="user.perm.admin == false">
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
    userNameClass() {
      return this.isNew ? "input input--block" :"input input--gray";
    },
    scopeClass() {
      if (this.$route.path == "/settings/global") {
        return "input input--block"
      }else{
        return this.isNew ? "input input--blocksub" :"input input--gray";
      }
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
        return `${baseClass} input--green`;
      }
    },
    passwordSubClass() {
      const baseClass = "input input--block";

      if ((this.user.password === "" || this.user.password === undefined) && this.passwordConf === "") {
        return baseClass;
      }
             
      if (this.user.password === this.passwordConf) {
        return `${baseClass} input--green`;
      }else{
        return `${baseClass} input--red`;
      }
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
