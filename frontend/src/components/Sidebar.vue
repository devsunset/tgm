<template>
  <nav :class="{ active }">
    <template v-if="isLogged">
      <router-link
        class="action"
        to="/files/"
        :aria-label="$t('sidebar.myFiles')"
        :title="$t('sidebar.myFiles')"
      >
        <i class="material-icons">folder</i>
        <span>{{ $t("sidebar.myFiles") }}</span>
      </router-link>

      <div v-if="user.perm.create && filesubmenu_visible" style="padding-left:20px">

        <button
          @click="openWebConsole() "
          class="action"
          :aria-label="$t('buttons.shell')"
          :title="$t('buttons.shell')"
        >
          <i class="material-icons">code</i>
          <span>{{ $t("sidebar.webconsole") }}</span>
        </button>

        <button
          @click="$store.commit('showHover', 'newDir')"
          class="action"
          :aria-label="$t('sidebar.newFolder')"
          :title="$t('sidebar.newFolder')"
        >
          <i class="material-icons">create_new_folder</i>
          <span>{{ $t("sidebar.newFolder") }}</span>
        </button>

        <button
          @click="$store.commit('showHover', 'newFile')"
          class="action"
          :aria-label="$t('sidebar.newFile')"
          :title="$t('sidebar.newFile')"
        >
          <i class="material-icons">note_add</i>
          <span>{{ $t("sidebar.newFile") }}</span>
        </button>

 
      </div>

      <div>
        <router-link
          v-if="user.perm.admin"
          class="action"
          to="/settings/users"
          :aria-label="$t('settings.userManagement')"
          :title="$t('settings.userManagement')"
        >
          <i class="material-icons">person</i>
          <span>{{ $t("settings.userManagement") }}</span>
        </router-link>

        <router-link
          v-if="user.perm.admin"
          class="action"
          to="/settings/groups"
          :aria-label="$t('sidebar.groupManagement')"
          :title="$t('sidebar.groupManagement')"
        >
          <i class="material-icons">group</i>
          <span>{{ $t("sidebar.groupManagement") }}</span>
        </router-link>

        <router-link
          class="action"
          to="/settings"
          :aria-label="$t('sidebar.settings')"
          :title="$t('sidebar.settings')"
        >
          <i class="material-icons">settings_applications</i>
          <span>{{ $t("sidebar.settings") }}</span>
        </router-link>

        <button
          v-if="authMethod == 'json'"
          @click="logout"
          class="action"
          id="logout"
          :aria-label="$t('sidebar.logout')"
          :title="$t('sidebar.logout')"
        >
          <i class="material-icons">exit_to_app</i>
          <span>{{ $t("sidebar.logout") }}</span>
        </button>
      </div>    
    </template>
    <template v-else>
      <router-link
        class="action"
        to="/login"
        :aria-label="$t('sidebar.login')"
        :title="$t('sidebar.login')"
      >
        <i class="material-icons">exit_to_app</i>
        <span>{{ $t("sidebar.login") }}</span>
      </router-link>

      <router-link
        v-if="signup"
        class="action"
        to="/login"
        :aria-label="$t('sidebar.signup')"
        :title="$t('sidebar.signup')"
      >
        <i class="material-icons">person_add</i>
        <span>{{ $t("sidebar.signup") }}</span>
      </router-link>
    </template>


    <p class="credits">
      <span>
          {{ $t("settings.username") }} : {{loginusername}}
      </span>  
      <span>
        <span v-if="disableExternal">TGM</span> 
        <a
          v-else
          rel="noopener noreferrer"
          target="_blank"
          href="https://www.teratec.co.kr/html/main.php"
          >TGM</a
        >
        <span>{{ version }}</span>
      </span>
      <span
        ><a @click="help">{{ $t("sidebar.help") }}</a></span
      >
    </p>
  </nav>
</template>

<script>
import { mapState, mapGetters } from "vuex";
import * as auth from "@/utils/auth";
import { webssh2port } from "@/utils/constants";
import {
  version,
  signup,
  disableExternal,
  noAuth,
  authMethod,
} from "@/utils/constants";

export default {
  name: "sidebar",
  data: () => ({ 
    filesubmenu_visible: true,
    loginusername : localStorage.getItem("username")
  }),
  computed: {
    ...mapState(["user"]),
    ...mapGetters(["isLogged"]),
    active() {
      return this.$store.state.show === "sidebar";
    },
    signup: () => signup,
    version: () => version,
    disableExternal: () => disableExternal,
    noAuth: () => noAuth,
    authMethod: () => authMethod,
    webssh2port: () => webssh2port,
  },
  methods: {
    help() {
      this.$store.commit("showHover", "help");
    },
    openWebConsole() {
       var x = localStorage.getItem("ssh") 
      if (x.split(",")[1] == "X") {
        alert(this.$t("settings.consolewarning1"));
        return;
      } 
      var admin_desc = "";
      if (x.split(",")[0] == "X") {
          if (this.user.perm.admin){
            admin_desc = "\n"+this.$t("settings.consolewarning2");
          }
      } else{
          if (this.user.perm.admin){
            admin_desc = "\n"+this.$t("settings.consolewarning3");
          }
      }

      alert(this.$t("settings.consolewarning")+admin_desc);

      var host = window.location.host;
      if (host.indexOf(":") > -1) {
        host = host.split(":")[0];
      }
      var currentTimeMillis = new Date().getTime();
      window.open(window.location.protocol+"//"+host +":"+this.webssh2port+"/ssh/host/"+host , "Web Console"+currentTimeMillis );  
      //  $store.commit('toggleShell')"

    },
    logout: auth.logout,
  },
  mounted() {
        if(this.$router.currentRoute.path.indexOf('/files/')>-1){
          this.filesubmenu_visible = true;
        }else{
          this.filesubmenu_visible = false;
        }
  },
  watch: {    
    $route: function () {
         if(this.$router.currentRoute.path.indexOf('/files/')>-1){
          this.filesubmenu_visible = true;
        }else{
          this.filesubmenu_visible = false;
        }
    },
  },
};
</script>
