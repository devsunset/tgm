<template>
  <div class="dashboard">
    <header-bar showMenu showLogo />

    <div id="nav">
      <div class="wrapper">

        <ul>
          <router-link to="/settings/profile" v-if="$route.path.indexOf('/settings/users') == -1 && $route.path.indexOf('/settings/groups') == -1"
            ><li :class="{ active: $route.path === '/settings/profile' }">
              {{ $t("settings.profileSettings") }}
            </li></router-link>
          <!-- 
            <router-link to="/settings/shares" v-if="user.perm.share"
            ><li :class="{ active: $route.path === '/settings/shares' }">
              {{ $t("settings.shareManagement") }}
            </li></router-link
          > 
          -->
          <router-link to="/settings/global" v-if="user.perm.admin && $route.path.indexOf('/settings/users') == -1 && $route.path.indexOf('/settings/groups') == -1"
            ><li :class="{ active: $route.path === '/settings/global' }">
              {{ $t("settings.globalSettings") }}
            </li></router-link>

          <router-link to="/settings/users" v-if="user.perm.admin  && $route.path.indexOf('/settings/users') > -1"
            ><li
              :class="{
                active:
                  $route.path.indexOf('/settings/users') > -1 || $route.name === 'User',
              }"
            >
              {{ $t("settings.userManagement") }}
            </li></router-link>


          <router-link to="/settings/groups" v-if="user.perm.admin  && $route.path.indexOf('/settings/groups') > -1"
            ><li
              :class="{
                active:
                  $route.path.indexOf('/settings/groups') > -1 || $route.name === 'Group',
              }"
            >
              {{ $t("sidebar.groupManagement") }}
          </li></router-link>

        </ul>

      </div>
    </div>

    <div v-if="loading">
      <h2 class="message delayed">
        <div class="spinner">
          <div class="bounce1"></div>
          <div class="bounce2"></div>
          <div class="bounce3"></div>
        </div>
        <span>{{ $t("files.loading") }}</span>
      </h2>
    </div>

    <router-view></router-view>

  </div>
</template>

<script>
import { mapState } from "vuex";

import HeaderBar from "@/components/header/HeaderBar";

export default {
  name: "settings",
  components: {
    HeaderBar,
  },
  computed: {
    ...mapState(["user", "loading"]),
  },
};
</script>
