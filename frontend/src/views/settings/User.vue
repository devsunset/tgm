<template>
  <errors v-if="error" :errorCode="error.message" />
  <div class="row" v-else-if="!loading">
    <div class="column">
      <form @submit="save" class="card">
        <div class="card-title">
          <h2 v-if="user.id === 0">{{ $t("settings.newUser") }}</h2>
          <h2 v-else>{{ $t("settings.user") }} {{ user.username }}</h2>
        </div>

        <div class="card-content">
          <user-form :user.sync="user" :isDefault="false" :isNew="isNew" ref="childform"/>
        </div>

        <div class="card-action">
          <button
            v-if="!isNew"
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
        if (this.isNew) {
          // USERNAME CHECK
          var regUserName = /^[a-z]+$/;
           if(false === regUserName.test(user.username)) {
            this.$showError('사용자 이름은 영문 소문자만 허용됩니다.');
            return;
          }

          // PASSWORD CHECK
          var pw = user.password
          var reg = /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$%^&*-]).{8,}$/;
          var hangulcheck = /[ㄱ-ㅎ|ㅏ-ㅣ|가-힣]/;
          if(false === reg.test(pw)) {
            this.$showError('비밀번호는 8자 이상이어야 하며, 숫자/대문자/소문자/특수문자를 모두 포함해야 합니다.');
            return;
          }else if(/(\w)\1\1\1/.test(pw)){
            this.$showError('같은 문자를 4번 이상 사용하실 수 없습니다.');
            return;
          }else if(pw.search(/\s/) != -1){
            this.$showError("비밀번호는 공백 없이 입력해주세요.");
            return;
          }else if(hangulcheck.test(pw)){
            this.$showError("비밀번호에 한글을 사용 할 수 없습니다."); 
            return;
          }else {
            if (pw !== this.$refs.childform.passwordConf || pw === "") {
              this.$showError("비밀번호 일치하지 않음"); 
              return;
            }
          }
                              
          const loc = await api.create(user);
          this.$router.push({ path: loc });
          this.$showSuccess(this.$t("settings.userCreated"));
        } else {
          // PASSWORD CHECK
          var pw = user.password
          var reg = /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$%^&*-]).{8,}$/;
          var hangulcheck = /[ㄱ-ㅎ|ㅏ-ㅣ|가-힣]/;
          if (pw !== ""){
            if(false === reg.test(pw)) {
              this.$showError('비밀번호는 8자 이상이어야 하며, 숫자/대문자/소문자/특수문자를 모두 포함해야 합니다.');
              return;
            }else if(/(\w)\1\1\1/.test(pw)){
              this.$showError('같은 문자를 4번 이상 사용하실 수 없습니다.');
              return;
            }else if(pw.search(/\s/) != -1){
              this.$showError("비밀번호는 공백 없이 입력해주세요.");
              return;
            }else if(hangulcheck.test(pw)){
              this.$showError("비밀번호에 한글을 사용 할 수 없습니다."); 
              return;
            }else {
              if (pw !== this.$refs.childform.passwordConf || pw === "") {
                this.$showError("비밀번호 일치하지 않음"); 
                return;
              }
            }
          }
          
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
