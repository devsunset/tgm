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
        v-model.trim="user.username"
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
        v-model.trim="user.password"
        name="password"
      />
      <input
          :class="passwordSubClass"
          type="password"
          :placeholder="$t('login.passwordConfirm')"
          v-model.trim="passwordConf"
          name="passwordConf"
      />
  </p>
    
<p v-if="!isDefault && user.perm.admin == false">
   <label for="shell">{{ $t("settings.shell") }} </label>
   <shells
      class="input input--block"
      id="shell"
      :shell.sync="user.shell"
    ></shells>
</p>

<p v-if="!isDefault && user.perm.admin == false">
  <label for="group">{{ $t("settings.group") }} </label>
    <multiselect v-model="value" placeholder="Select group"   label="name" track-by="code" :options="options" :multiple="true" :taggable="false"></multiselect>
 </p>

 <p v-if="!isDefault && user.perm.admin == false">
  <label>{{ $t("settings.accountexpireday") }} </label>
      <date-picker v-model="user.expireDay"  :disabled-date="disabledBeforeTodayAndAfterAWeek" value-type="format" format="YYYY-MM-DD" placeholder="Select date"></date-picker>
 </p>

  <p v-if="!isDefault && user.perm.admin == false">
  <label>{{ $t("settings.passwordexpireday") }} </label>
  <input type="number" class="input" style="width:212px" v-model="user.passwordExpireDay" min="-1" max="365"   pattern="^[0-9]" onkeypress="return (event.charCode == 8 || event.charCode == 0 || event.charCode == 13) ? null : event.charCode >= 48 && event.charCode <= 57"> {{ $t("time.days") }}<br>
 </p>

  <p v-if="!isDefault && user.perm.admin == false">
  <label>{{ $t("settings.passwordexpirewarningday") }} </label>
  <input type="number" class="input" style="width:212px" v-model="user.passwordExpireWarningDay" min="-1" max="30"     pattern="^[0-9]" onkeypress="return (event.charCode == 8 || event.charCode == 0 || event.charCode == 13) ? null : event.charCode >= 48 && event.charCode <= 57"> {{ $t("time.days") }}<br>
 </p>

<p v-if="!isNew && !isDefault &&  user.perm.admin == false">
   <input
        type="checkbox"
        :disabled="user.perm.admin"
        v-model="user.lockAccount"
      />
      {{ $t("settings.accountlock") }}
</p>

<p>
      <label for="scope">{{ $t("settings.scope") }}</label>
      <input
        :class="scopeClass"
        type="text"
        v-model="user.scope"
        id="scope"
        :disabled="!this.isNew && $route.path != '/settings/global'"
        style="display:inline-block"
      />
</p>
<p v-if="this.isNew">
    Home path -> {{user.scope}}<span v-if="user.username">/</span>{{user.username}}
</p>
<p>
     <progress-bar
        :options="baroptions"
        :value="homeuse"
      />
</p>
<p>
    [ /home ] -> Size :  {{homesize}}, Used :  {{homeused}}, Avail :  {{homeavail}}
</p>

<br>
<!-- 
  <br>
  <p v-if="!isDefault && user.perm.admin == false">
    <label>[[[Quota - 요구사항 Fix 후 진행 예정]]]</label>
  </p>
  <br> 
-->

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
    <p v-if="!isDefault" v-show="false">
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
import { users as api } from "@/api";

import Languages from "./Languages";
import Shells from "./Shells";
 import Multiselect from 'vue-multiselect'
import Rules from "./Rules";
import Permissions from "./Permissions";
import Commands from "./Commands";
import { enableExec } from "@/utils/constants";
import { enableCmdLimit } from "@/utils/constants";
import DatePicker from 'vue2-datepicker';
import 'vue2-datepicker/index.css';
 import ProgressBar from 'vuejs-progress-bar'

export default {
  name: "user",
  data: function () {
    return {      
      passwordConf: "",
      value: [],
      options: [],
      groupsvalue: [],
      baroptions: {
        text: {
          color: '#FFFFFF',
          shadowEnable: true,
          shadowColor: '#000000',
          fontSize: 14,
          fontFamily: 'Helvetica',
          dynamicPosition: false,
          hideText: false
        },
        progress: {
          color: '#2dbd2d',
          backgroundColor: '#333333',
          inverted: false
        },
        layout: {
          height: 35,
          width: 212,
          verticalTextAlign: 61,
          horizontalTextAlign: 43,
          zeroOffset: 0,
          strokeWidth: 30,
          progressPadding: 0,
          type: 'line'
        }
    },
    homesize: "0G",
    homeused: "0G",
    homeavail: "0G",
    homeuse: 0,
    };
  },
  components: {
    Permissions,
    Languages,
    Shells,
    Multiselect,
    DatePicker,
    Rules,
    Commands,
    ProgressBar,
  },
  props: ["user", "isNew", "isDefault"],
  methods : {
    disabledBeforeTodayAndAfterAWeek(date) {
      const today = new Date();
      today.setHours(0, 0, 0, 0);
      //return date < today || date > new Date(today.getTime() + 7 * 24 * 3600 * 1000);
       return date < today;
    },
  },
  async created() {
    try {
          this.groupsvalue = await api.getGroups();
          for (const [key, value] of Object.entries(this.groupsvalue)) {
            this.options.push({
              code: key,
              name: value,
            });
          }

        if (this.user.group != null && this.user.group != "") {
          var sValue = this.user.group.split(",");
          for (var i = 0; i < sValue.length; i++) {
            this.value.push({
                code: sValue[i],
                name: sValue[i],
              });
          } 
        }

        try {
              var home = await api.getHomeInfo("/home");
              var shome = home.split(",");
              this.homesize = shome[0];
              this.homeused = shome[1];
              this.homeavail = shome[2];
              this.homeuse = shome[3].replace("%", "");
        } catch (e) {
          console.log(e);
        }

    } catch (e) {
      this.error = e;
    }
  },
  computed: {
    passwordPlaceholder() {
      return this.isNew ? this.$t("settings.password") : this.$t("settings.avoidChanges");
    },
    userNameClass() {
      return this.isNew ? "input input--block" :"input input--gray";
    },
    scopeClass() {
      if (this.$route.path == "/settings/global") {
        return "input input--block"
      }else{
        return this.isNew ? "input input--block" :"input input--gray";
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
 mounted: function () {
    if (this.user.expireDay == null || this.user.expireDay == undefined || this.user.expireDay == "") {
        let date = new Date();
        let month = date.getMonth() + 1;
        let day = date.getDate();
        month = month >= 10 ? month : '0' + month;
        day = day >= 10 ? day : '0' + day;
        this.user.expireDay = date.getFullYear()+1 + '-' + month + '-' + day ;
    }
    if (this.user.passwordExpireDay == null || this.user.passwordExpireDay == undefined || this.user.passwordExpireDay == "") {
        this.user.passwordExpireDay = 90;
    }
    if (this.user.passwordExpireWarningDay== null || this.user.passwordExpireWarningDay == undefined || this.user.passwordExpireWarningDay == "") {
          this.user.passwordExpireWarningDay = 7;
    }
  },
  watch: {
    "user.perm.admin": function () {
      if (!this.user.perm.admin) return;
      this.user.lockPassword = false;
    },
    "user.passwordExpireDay": function () {
       if (this.user.passwordExpireDay > 365) {
         this.user.passwordExpireDay = 365;
        }
       if (this.user.passwordExpireDay < -1) {
          this.user.passwordExpireDay = -1;
        }
    },
      "user.passwordExpireWarningDay": function () {
       if (this.user.passwordExpireWarningDay > 30) {
          this.user.passwordExpireWarningDay = 30;
        }
        if (this.user.passwordExpireWarningDay < -1) {
          this.user.passwordExpireWarningDay = -1;
        }
    },
  },
};
</script>

<!-- New step!
     Add Multiselect CSS. Can be added as a static asset or inside a component. -->
<style src="vue-multiselect/dist/vue-multiselect.min.css"></style>

<style>
.mx-input {
  background-color :  #3A4147 !important;
  color: rgba(255, 255, 255, 0.87) !important;
}

.mx-calendar.mx-calendar-panel-date {
  background-color :  #3A4147 !important;
  color: rgba(255, 255, 255, 0.87) !important;
}

.multiselect__content-wrapper {
  background-color :  #3A4147 !important;
  color: rgba(255, 255, 255, 0.87) !important;
}

.multiselect__tags {
  background-color :  #3A4147 !important;
  color: rgba(255, 255, 255, 0.87) !important;
}
</style>
