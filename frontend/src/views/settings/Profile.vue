<template>
  <div class="row">
    <div class="column">
      <form class="card" @submit="updateSettings">
        <div class="card-title">
          <h2>{{ $t("settings.profileSettings") }}</h2>
        </div>

        <div class="card-content">
          <p>
            <input type="checkbox" v-model="hideDotfiles" />
            {{ $t("settings.hideDotfiles") }}
          </p>
          <p>
            <input type="checkbox" v-model="singleClick" />
            {{ $t("settings.singleClick") }}
          </p>
          <h3>{{ $t("settings.language") }}</h3>
          <languages
            class="input input--block"
            :locale.sync="locale"
          ></languages>
        </div>

        <div class="card-action">
          <input
            class="button button--flat"
            type="submit"
            :value="$t('buttons.update')"
          />
        </div>
      </form>
    </div>

    <div class="column">
      <form class="card" v-if="!user.lockPassword" @submit="updatePassword">
        <div class="card-title">
          <h2>{{ $t("settings.changePassword") }}</h2>
        </div>

        <div class="card-content">
          <input
            :class="passwordClass"
            type="password"
            :placeholder="$t('settings.newPassword')"
            v-model="password"
            name="password"
          />
          <input
            :class="passwordSubClass"
            type="password"
            :placeholder="$t('settings.newPasswordConfirm')"
            v-model="passwordConf"
            name="passwordConf"
          />
        </div>

        <div class="card-action">
          <input
            class="button button--flat"
            type="submit"
            :value="$t('buttons.update')"
          />
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { mapState, mapMutations } from "vuex";
import { users as api } from "@/api";
import Languages from "@/components/settings/Languages";

export default {
  name: "settings",
  components: {
    Languages,
  },
  data: function () {
    return {
      password: "",
      passwordConf: "",
      hideDotfiles: false,
      singleClick: false,
      locale: "",
    };
  },
  computed: {
    ...mapState(["user"]),
    passwordClass() {
      const baseClass = "input input--block";

      if (this.password === "" && this.passwordConf === "") {
        return baseClass;
      }

      var pw = this.password
              
      var reg = /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$%^&*-]).{8,}$/;
      var hangulcheck = /[ㄱ-ㅎ|ㅏ-ㅣ|가-힣]/;
      
      if(false === reg.test(pw)) {
        return `${baseClass} input--red`;
      }else if(/(\w)\1\1\1/.test(pw)){
       return `${baseClass} input--red`;
      }else if(pw.search(/\s/) != -1){
        return `${baseClass} input--red`;
      }else if(hangulcheck.test(pw)){
       return `${baseClass} input--red`;
      }else {
        return `${baseClass} input--green`;
      }
    },
    passwordSubClass() {
      const baseClass = "input input--block";
      if (this.password === ""  && this.passwordConf === "") {
        return baseClass;
      }
             
      if (this.password === this.passwordConf) {
        return `${baseClass} input--green`;
      }else{
        return `${baseClass} input--red`;
      }
    },
  },
  created() {
    this.setLoading(false);
    this.locale = this.user.locale;
    this.hideDotfiles = this.user.hideDotfiles;
    this.singleClick = this.user.singleClick;
  },
  methods: {
    ...mapMutations(["updateUser", "setLoading"]),
    async updatePassword(event) {

      event.preventDefault();

      var pw = this.password
              
      var reg = /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$%^&*-]).{8,}$/;
      var hangulcheck = /[ㄱ-ㅎ|ㅏ-ㅣ|가-힣]/;
      
      if(false === reg.test(pw)) {
        this.$showError(this.$t("login.passwordrule1"));
        return;
      }else if(/(\w)\1\1\1/.test(pw)){
        this.$showError(this.$t("login.passwordrule2"));
        return;
      }else if(pw.search(/\s/) != -1){
        this.$showError(this.$t("login.passwordrule3"));
        return;
      }else if(hangulcheck.test(pw)){
        this.$showError(this.$t("login.passwordrule4"));
        return;
      }else {
        if (this.password !== this.passwordConf || this.password === "") {
          this.$showError(this.$t("login.passwordrule5")); 
          return;
        }
        try {
          const data = { id: this.user.id, password: this.password };
          await api.update(data, ["password"]);
          this.updateUser(data);
          this.$showSuccess(this.$t("settings.passwordUpdated"));
        } catch (e) {
          this.$showError(e);
        }
      }
    },
    async updateSettings(event) {
      event.preventDefault();

      try {
        const data = {
          id: this.user.id,
          locale: this.locale,
          hideDotfiles: this.hideDotfiles,
          singleClick: this.singleClick,
        };
        await api.update(data, ["locale", "hideDotfiles", "singleClick"]);
        this.updateUser(data);
        this.$showSuccess(this.$t("settings.settingsUpdated"));
      } catch (e) {
        this.$showError(e);
      }
    },
  },
};
</script>
