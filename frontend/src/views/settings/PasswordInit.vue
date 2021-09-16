<template>
  <errors v-if="error" :errorCode="error.message" />
  <div class="row" v-else-if="!loading">
    <div class="column">
      <div class="card">
        <div class="card-title">
          <h2>{{ $t("sidebar.requestpasswordreq") }} - 신규개발 필요(사용자관리 메뉴로 복사만 한 상태)</h2>
        </div>

        <div class="card-content full">
          <table>
            <tr>
              <th>{{ $t("settings.username") }}</th>
              <th>Request Date</th>
              <th>Status</th>
              <th></th>
            </tr>

            <tr v-for="user in users" :key="user.id">
              <td>{{ user.username }}</td>
              <td>2021-09-18</td>
              <td>Request</td>
              <td class="small">
                <router-link :to="'/settings/passwordinit/' + user.id"
                  ><i class="material-icons">mode_edit</i></router-link
                >
              </td>
            </tr>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState, mapMutations } from "vuex";
import { users as api } from "@/api";
import Errors from "@/views/Errors";

export default {
  name: "users",
  components: {
    Errors,
  },
  data: function () {
    return {
      error: null,
      users: [],
    };
  },
  async created() {
    this.setLoading(true);

    try {
      this.users = await api.getAll();
    } catch (e) {
      this.error = e;
    } finally {
      this.setLoading(false);
    }
  },
  computed: {
    ...mapState(["loading"]),
  },
  methods: {
    ...mapMutations(["setLoading"]),
  },
};
</script>
