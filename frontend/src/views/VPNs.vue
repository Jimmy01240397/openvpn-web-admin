<template>
  <div>
    <v-expansion-panels class="mb-6">
      <v-expansion-panel
        v-for="(vpndata, i) in vpndatas"
        :key="i"
      >
        <v-expansion-panel-title expand-icon="mdi-plus" collapse-icon="mdi-minus">
          <v-row no-gutters>
            <v-col cols="4" class="d-flex justify-start">
              {{ vpndata.VPNname }}
            </v-col>
            <v-col
              cols="8"
              class="text--secondary"
            >
              Status: {{ vpndata.Active ? "Active" : "Stop" }}
            </v-col>
          </v-row>
        </v-expansion-panel-title>
        <v-expansion-panel-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              variant="text"
              color="red"
            >
              Delete

              <v-dialog
                v-model="dialog"
                activator="parent"
                width="auto"
              >
                <v-card>
                  <v-card-text>
                    Do you realy want to delete {{ vpndata.VPNname }}?
                  </v-card-text>
                  <v-card-actions>
                    <v-btn color="gray" @click="dialog = false">Cancel</v-btn>
                    <v-btn color="red" @click="deletevpn(vpndata)">Delete</v-btn>
                  </v-card-actions>
                </v-card>
              </v-dialog>
            </v-btn>
            <v-btn
              variant="text"
              color="blue"
              @click="startstopvpn(vpndata)"
            >
              {{ vpndata.Active ? "Stop" : "Start" }}
            </v-btn>
          </v-card-actions>
        </v-expansion-panel-text>
      </v-expansion-panel>
    </v-expansion-panels>
    <h2>Add New VPN</h2>
    <div class="message">
      <v-alert
        :color='messagecolor'
        :type='messagetype'
        :model-value='showmessage'
      >
        {{ message }}
      </v-alert>
    </div>
    <form class='addvpn-form'>
      <v-text-field
        label='VPNname'
        v-model='vpnname'
        @keyup.enter='addvpn'
      ></v-text-field>
      <v-btn
        @click='addvpn'
      >
        Add
      </v-btn>
    </form>

  </div>
</template>

<script>
import vpnService from '@/services/vpn';
export default {
  name: 'VPNS',
  data: () => ({
    vpnname: '',
    //password: '',
    //checkpassword: '',
    messagecolor: 'red',
    messagetype: 'error',
    showmessage: false,
    message: "",
    dialog: false,
    vpndatas: [],
  }),
  methods: {
    addvpn: async function() {
      const result = await vpnService.addvpn(this.vpnname, true);
      switch (result.status) {
        case 200:
          this.vpnname = '';
          this.messagecolor = 'green';
          this.messagetype = 'success'
          this.showmessage = true;
          this.message = "Add vpn success!";
          this.vpndatas = (await vpnService.getvpns()).data;
          break;
        default:
          this.messagecolor = 'red';
          this.messagetype = 'error';
          this.showmessage = true;
          this.message = result.data;
          break;
      }
    },
    deletevpn: async function(vpndata) {
      const result = await vpnService.deletevpn(vpndata.VPNname);
      switch (result.status) {
        case 200:
          this.vpndatas = (await vpnService.getvpns()).data;
          break;
        default:
          this.messagecolor = 'red';
          this.messagetype = 'error';
          this.showmessage = true;
          this.message = result.data;
          break;
      }
      this.dialog = false;
    },
    startstopvpn: async function(vpndata) {
      const result = await vpnService.startstopvpn(vpndata.VPNname, !vpndata.Active);
      switch (result.status) {
        case 200:
          this.vpndatas = (await vpnService.getvpns()).data;
          break;
        default:
          this.messagecolor = 'red';
          this.messagetype = 'error';
          this.showmessage = true;
          this.message = result.data;
          break;
      }
    },
  },
  async beforeMount() {
    this.vpndatas = (await vpnService.getvpns()).data;
  }
}
</script>
