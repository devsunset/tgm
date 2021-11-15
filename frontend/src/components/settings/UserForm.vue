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
   <shells
      class="input input--block"
      id="shell"
      :shell.sync="user.shell"
    ></shells>
</p>

<p v-if="user.perm.admin == false">
  <label for="group">Group - 작업중</label>
   <groups
      class="input input--block"
      id="group"
      :group.sync="user.group"
    ></groups>
 </p>

 <p v-if="user.perm.admin == false">
  <label>계정 유효 일자 - 작업중</label>
  <input type="date" id="start" name="trip-start"
       value="2021-11-15"   min="2021-11-15" max="2022-12-31">
 </p>

  <p v-if="user.perm.admin == false">
  <label>암호 기간 만료일 - 작업중</label>
  <input type="text" @input="bindNumber1" :value="number1" maxlength="4" style="width: 50px;"> 일<br>
 </p>

  <p v-if="user.perm.admin == false">
  <label>암호 변경 경고일 - 작업중</label>
  <input type="text" @input="bindNumber2" :value="number2" maxlength="4" style="width: 50px;"> 일<br>
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
  <label>Quota - 요구사항 정의 fix 후 진행 예정</label>
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
import Shells from "./Shells";
import Groups from "./Groups";
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
      number1: '',
      number2: '',
    };
  },
  components: {
    Permissions,
    Languages,
    Shells,
    Groups,
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
  methods:{
    bindNumber1(event){
      this.number1 = event.target.value;
    },
    bindNumber2(event){
      this.number2= event.target.value;
    },
  },
  watch: {
    "user.perm.admin": function () {
      if (!this.user.perm.admin) return;
      this.user.lockPassword = false;
    },
   number1(val){
      const reg = /[ㄱ-ㅎ|ㅏ-ㅣ|가-힣|a-z]/;

      // 한글, 영문 체크
      if(reg.exec(val)!==null) this.number1 = val.replace(/[^0-9]/g,'');

      this.number1 = val.replace('-','');

      // .... 만 입력하게 될 경우 체크
      if(isNaN(parseFloat(val))) this.number1 = '';
    },
    number2(val){
      const reg = /[ㄱ-ㅎ|ㅏ-ㅣ|가-힣|a-z]/;

      // 한글, 영문 체크
      if(reg.exec(val)!==null) this.number2 = val.replace(/[^0-9]/g,'');

      // .... 만 입력하게 될 경우 체크
      if(isNaN(parseFloat(val))) this.number2 = '';
    },
  },
};
</script>
