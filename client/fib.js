new Vue({
  el: '#app',
  data: {
    number: 0,
    data: undefined,
    error: undefined
  },
  methods: {
    getFib() {
      this.number = this.number ? this.number: 0;
      axios.get(`/api/fib/${this.number}`).then((response) => {
        this.data = response.data;
        this.error = undefined;
      })
      .catch((err) => {
        this.data = undefined;
        this.error = err.response.data;
      });
    }
  }
})
