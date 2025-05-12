<template>
  <div>
    <h2 class="mx-5 my-5">GoNews - агрегатор новостей.</h2>
    <div v-for="post in news" :key="post.ID">
      <v-card elevation="10" outlined class="mx-5 my-5">
        <v-card-title>
          <a :href="post.Link" target="_blank"> {{ post.Title }} </a>
        </v-card-title>
        <v-card-text>
          {{ post.Content }}
          <v-card-subtitle>
            {{ new Date(post.PubTime * 1000) }}
          </v-card-subtitle>
        </v-card-text>
      </v-card>
    </div>
  </div>
</template>

<script>
export default {
  name: "News",
  data() {
    return {
      news: [],
    };
  },
  mounted() {
    let url = "http://" + window.location.hostname + "/news/40";
    fetch(url)
      .then((response) => response.json())
      .then((data) => (this.news = data));
  },
};
</script>

<style scoped>
</style>
