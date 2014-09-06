var getWordsAndDraw = function () {
  $.ajax({
    type: "GET",
    cache: false,
    url: "/words",

    success: function (dataobj, state, jqxhr) {
      var attach = $('#results');
      attach.html("");
      var lines = '<tr >' +
        '<th>Word</th>' +
        '<th>Pos</th>' +
        '<th>Definition</th>' +
        '<th>Language</th>' +
        '<th>Example</th>' +
        '<th>Ethymology</th>' +
        '<th>Delete</th>' +
        '</tr>';

      data = JSON.parse(dataobj);

      if (!data.Result) {
        $("#status-message")
          .html("Your Dictionary is Emtpy, click <a href='/add'>Here</a> to begin feeding it !")
          .addClass('alert-info');
        return;
      } else {
        $("#status-message").text("Words in your LittleDictionary");
      }

      data.Words.forEach(function (o) {
        if (o.Word == "" || !o.Word) {
          return;
        }

        var example = o.Example || "";
        var ethymology = o.Ethymology || "";

        var line = '<tr id="tr_' + o.Id + '">' +
          '<td id="' + o.Id + '_Word" class="editable word">' + o.Word + '</td>' +
          '<td id="' + o.Id + '_Pos" class="pos_editable">' + o.Pos + '</td>' +
          '<td id="' + o.Id + '_Definition" class="editable">' + o.Definition + '</td>' +
          '<td id="' + o.Id + '_Language" class="editable">' + o.Language + '</td>' +
          '<td id="' + o.Id + '_Example" class="editable">' + example + '</td>' +
          '<td id="' + o.Id + '_Ethymology" class="editable">' + ethymology + '</td>' +
          '<td id="' + o.Id + '" class="delete">x</td>' +
          '</tr>';
        lines += line;
      });
      attach.html(lines);
      bindEdit();
      bindPosEdit();
      bindDelete();

    },
    error: function (jqxhr, state, err) {
      $("#status-message").text("Internal Error")
      $('#results').innerHTML = ""
    }
  });
}

var bindDelete = function () {
  $('.delete').click(function (event) {
    var p = $(event.target).parent();
    wordId = event.target.id;
    $.ajax({
      type: "POST",
      cache: false,
      url: "/word/delete",
      data: JSON.stringify({
        Id: wordId
      }),
      success: function (dataobj, state, jqxhr) {
        p.remove();
      },
      error: function (jqxhr, state, err) {
        $("#status-message").text("Internal Error")
        $('#results').innerHTML = ""
      }
    });
  });
}

var bindPosEdit = function () {

  $('.pos_editable').editable(function (value, settings) {
    var s = $(this).attr('id').split('_');
    var wordId = s[0];
    var field = s[1];
    console.log(value, wordId, field);
    console.log(this);

    var data = {
      "Id": wordId
    }

    data[field] = value;

    console.log(data);

    $.ajax({
      type: "POST",
      cache: false,
      url: "/word/update",
      data: JSON.stringify(data),
      success: function (dataobj, state, jqxhr) {},
      error: function (jqxhr, state, err) {
        $("#status-message").text("Internal Error")
        $('#results').innerHTML = ""
      }
    });

    return value;
  }, {
    type: 'select',
    data: {
      Noun: "Noun",
      Pronoun: "Pronoun",
      Verb: "Verb",
      Adjective: "Adjective",
      Adverb: "Adverb",
      Preposition: "Preposition",
      Conjunction: "Conjunction",
      Interjection: "Interjection",
      Expression: "Expression"
    },

    submit: 'Ok',
    cancel: 'Cancel'
  })
}

var bindEdit = function () {

  $('.editable').editable(function (value, settings) {
    var s = $(this).attr('id').split('_');
    var wordId = s[0];
    var field = s[1];
    console.log(value, wordId, field);
    console.log(this);

    var data = {
      "Id": wordId
    }

    data[field] = value;

    console.log(data);

    $.ajax({
      type: "POST",
      cache: false,
      url: "/word/update",
      data: JSON.stringify(data),
      success: function (dataobj, state, jqxhr) {},
      error: function (jqxhr, state, err) {
        $("#status-message").text("Internal Error")
        $('#results').innerHTML = ""
      }
    });
    return value;
  }, {
    submit: 'Ok',
    cancel: 'Cancel'
  })
}