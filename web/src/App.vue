<template>
  <v-app>
    <v-main>
      <News />
    </v-main>
  </v-app>
</template>

<script>
import News from "./components/News";

export default {
  name: "App",
  components: {
    News,
  },

  data: () => ({
    //
  }),
};
</script>
