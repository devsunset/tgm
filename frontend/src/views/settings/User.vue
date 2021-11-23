<template>
  <errors v-if="error" :errorCode="error.message" />
  <div class="row" v-else-if="!loading">
    <div class="column" style="width:100%">
      <form @submit="save" class="card">
        <div class="card-title">
          <h2 v-if="user.id === 0">{{ $t("settings.newUser") }}  - [개발 진행 중]</h2>
          <h2 v-else>{{ $t("settings.user") }} {{ user.username }} - [개발 진행 중]</h2>
        </div>

        <div class="card-content">
          <user-form :user.sync="user" :isDefault="false" :isNew="isNew" ref="childform"/>
        </div>

        <div class="card-action">
          <button
            v-if="!isNew && user.username !='admin'"
            @click.prevent="deletePrompt"
            type="button"
            class="button button--flat button--red"
            :aria-label="$t('buttons.delete')"
            :title="$t('buttons.delete')"
          >
            {{ $t("buttons.delete") }}
          </button>
          <input
            class="button button--flat"
            type="submit"
            :value="$t('buttons.save')"
          />
        </div>
      </form>
    </div>

    <div v-if="$store.state.show === 'deleteUser'" class="card floating">
      <div class="card-content">
        <p>Are you sure you want to delete this user?</p>
      </div>

      <div class="card-action">
        <button
          class="button button--flat button--grey"
          @click="closeHovers"
          v-focus
          :aria-label="$t('buttons.cancel')"
          :title="$t('buttons.cancel')"
        >
          {{ $t("buttons.cancel") }}
        </button>
        <button class="button button--flat" @click="deleteUser">
          {{ $t("buttons.delete") }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState, mapMutations } from "vuex";
import { users as api, settings } from "@/api";
import UserForm from "@/components/settings/UserForm";
import Errors from "@/views/Errors";
import deepClone from "lodash.clonedeep";

export default {
  name: "user",
  components: {
    UserForm,
    Errors,
  },
  data: () => {
    return {
      error: null,
      originalUser: null,
      user: {},
    };
  },
  created() {
    this.fetchData();
  },
  computed: {
    isNew() {
      return this.$route.path === "/settings/users/new";
    },
    ...mapState(["loading"]),
  },
  watch: {
    $route: "fetchData",
    "user.perm.admin": function () {
      if (!this.user.perm.admin) return;
      this.user.lockPassword = false;
    },
  },
  methods: {
    ...mapMutations(["closeHovers", "showHover", "setUser", "setLoading"]),
    async fetchData() {
      this.setLoading(true);

      try {
        if (this.isNew) {
          let { defaults } = await settings.get();
          this.user = {
            ...defaults,
            username: "",
            password: "",
            shell: "/bin/bash",
            group : "",
            expireDay : "",
            passwordExpireDay : "",
            passwordExpireWarningDay : "",
            lockAccount : false,
            rules: [],
            lockPassword: false,
            id: 0,
          };
        } else {
          const id = this.$route.params.pathMatch;
          this.user = { ...(await api.get(id)) };
        }
      } catch (e) {
        this.error = e;
      } finally {
        this.setLoading(false);
      }
    },
    deletePrompt() {
      this.showHover("deleteUser");
    },
    async deleteUser(event) {
      event.preventDefault();

      try {
        await api.remove(this.user.id);
        this.$router.push({ path: "/settings/users" });
        this.$showSuccess(this.$t("settings.userDeleted"));
      } catch (e) {
        e.message === "403"
          ? this.$showError(this.$t("errors.forbidden"), false)
          : this.$showError(e);
      }
    },
    async save(event) {
      event.preventDefault();
      let user = {
        ...this.originalUser,
        ...this.user,
      };

      try {
          // PASSWORD CHECK
        var pw = user.password
        var reg = /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$%^&*-]).{8,}$/;
        var hangulcheck = /[ㄱ-ㅎ|ㅏ-ㅣ|가-힣]/;

        if (this.isNew) {
          // USERNAME CHECK
          var regUserName = /^[A-Za-z0-9+]*$/;
           if(false === regUserName.test(user.username)) {
            this.$showError(this.$t("settings.usernamerule"));
            return;
          }

          var blank_pattern = /^\s+|\s+$/g;
          if( user.username == '' || user.username == null || user.username.replace( blank_pattern, '' ) == "" ){
            this.$showError(this.$t("settings.inputusername"));
            return;
          }


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
            if (pw !== this.$refs.childform.passwordConf || pw === "") {
              this.$showError(this.$t("login.passwordrule5")); 
              return;
            }
          }
        } else {

          if (pw !== ""){
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
              if (pw !== this.$refs.childform.passwordConf || pw === "") {
                this.$showError(this.$t("login.passwordrule5")); 
                return;
              }
            }
          }
        }

        alert("[[[ To-Do ]]] ---> Check Validation");
        alert("shell: "+this.user.shell)
        alert("group: "+this.user.group)
        alert("expireDay: "+this.user.expireDay)
        alert("passwordExpireDay: "+this.user.passwordExpireDay)
        alert("passwordExpireWarningDay: "+this.user.passwordExpireWarningDay)
        alert("lockAccount: "+this.user.lockAccount)
        alert("home: "+this.user.scope)
        alert(JSON.stringify(user))

        if (this.isNew) {
          const loc = await api.create(user);
          this.$router.push({ path: loc });
          this.$showSuccess(this.$t("settings.userCreated"));
        } else {
          await api.update(user);
          if (user.id === this.$store.state.user.id) {
            this.setUser({ ...deepClone(user) });
          }
          this.$showSuccess(this.$t("settings.userUpdated"));
        }
      } catch (e) {
        this.$showError(e);
      }
    },
  },
};
</script>
