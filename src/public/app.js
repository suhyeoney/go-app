new Vue({
    el: '#app',

    data: {
        websocket: null,
        newMessage: '',
        chatContent: '',
        email: null,
        username: null,
        joined: false
    },
    created: function() {
        var self = this;
        this.websocket = new WebSocket('ws://' + window.location.host + '/ws');
        this.websocket.addEventListener('message', function(e) {
            var message = JSON.parse(e.data);
            self.chatContent += '<div class="chip">'
                + '<img src="' + self.gravatarURL(message.email) + '">'
                + message.username
                + '</div>'
                + message.message + '<br/>';

            var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight;
        });
    },
    methods: {
        send: function () {
            if (this.newMessage != '') {
                this.websocket.send(
                    JSON.stringify({
                            email: this.email,
                            username: this.username,
                            message: $('<p>').html(this.newMessage).text()
                        }
                    ));
                this.newMessage = '';
            }
        },
        join: function () {
            if (!this.email) {
                Materialize.toast('You must enter an email', 2000);
                return;
            }
            if (!this.username) {
                Materialize.toast('You must choose a username', 2000);
                return;
            }
            this.email = $('<p>').html(this.email).text();
            this.username = $('<p>').html(this.username).text();
            this.joined = true;
        },
        gravatarURL: function(email) {
            return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
        }
    }
});