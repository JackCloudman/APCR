(function () {
  var Message;
  Message = function (arg) {
      this.text = arg.text, this.message_side = arg.message_side;
      this.draw = function (_this) {
          return function () {
              var $message;
              $message = $($('.message_template').clone().html());
              $message.addClass(_this.message_side).find('.text').html(_this.text);
              $('.messages').append($message);
              return setTimeout(function () {
                  return $message.addClass('appeared');
              }, 0);
          };
      }(this);
      return this;
  };
  $(function () {
      var getMessageText, message_side, sendMessage;

// Call Python function, and pass explicit callback function
      getMessageText = function () {
          var $message_input;
          $message_input = $('.message_input');
          return $message_input.val();
      };
      sendMessage = function (text,side) {
        if(side=="right"){
          eel.sendCommand(text,"message");  // Call a Python function
        }
          var $messages, message;
          if (text.trim() === '') {
              return;
          }
          $('.message_input').val('');
          $messages = $('.messages');
          message = new Message({
              text: text,
              message_side: side
          });
          message.draw();
          return $messages.animate({ scrollTop: $messages.prop('scrollHeight') }, 300);
      };
      eel.expose(recvMessage);               // Expose this function to Python
       function recvMessage(text) {
           sendMessage(text,"left")
       }
      $('.send_message').click(function (e) {
          return sendMessage(getMessageText(),"right");
      });
      $('.message_input').keyup(function (e) {
          if (e.which === 13) {
              return sendMessage(getMessageText(),"right");
          }
      });

  });
}.call(this));
