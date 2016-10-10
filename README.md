# api.livereload.com
LiveReload server-side API, the new version (in Go)


## Paddle endpoint

Example POST request (`application/x-www-form-urlencoded`):

    token = ...
    name = Test User
    email = 57fb78dc8c9f6@blackhole.io
    quantity = 1
    txn = e5043ee75a37b585449e4ccfff0195ab
    message = 55f3b051cc88a8b5a55ed815e046aa14
    passthrough = Example passthrough
    p_earnings = {"128":"9.0000"}
    p_paddle_fee = 1
    p_price = 10
    p_quantity = 1
    p_signature = ...
    p_coupon = 
    p_currency = USD
    p_order_id = 922909
    p_product_id = 489469
    p_sale_gross = 10
    p_tax_amount = 0
    p_country = US
    p_coupon_savings = 0
    p_used_price_override = 1

## TODO: beta signup

Email:

    Subject: Beta: NAME <EMAIL>

    ABOUT
    ---
    LiveReload APP_VERSION (APP_PLATFORM)
    NAME <EMAIL>

Route:

    var NEW_BETA_SIGNUP_EMAIL = require('fs').readFileSync(__dirname + '/../views/emails/beta-signup.txt', 'utf8');
    app.post('/api/v1/beta-signup/', wrap((ctx, req, res) => {
    var params = {
    appPlatform: (req.param('appPlatform') || '').trim(),
    appVersion: (req.param('appVersion') || '').trim(),
    name: (req.param('name') || '').trim(),
    email: (req.param('email') || '').trim(),
    about: (req.param('about') || '').trim()
    };

    if (params.name == '')
    return res.send(400, { error: 'EREQ', message: 'Missing name parameter' });
    if (params.email == '')
    return res.send(400, { error: 'EREQ', message: 'Missing email parameter' });
    if (params.appPlatform == '')
    return res.send(400, { error: 'EREQ', message: 'Missing appPlatform parameter' });
    if (!params.appPlatform.match(/^mac|windows$/))
    return res.send(400, { error: 'EREQ', message: 'Invalid appPlatform parameter' });
    if (params.appVersion == '')
    return res.send(400, { error: 'EREQ', message: 'Missing appVersion parameter' });
    if (!params.appVersion.match(/^\d+\.\d+\.\d+$/))
    return res.send(400, { error: 'EREQ', message: 'Invalid appVersion parameter' });

    model.recordBetaSignup(ctx, params).then((signUpRecord) => {
    res.send({ ok: 'ok' });

    var info = { RECORD_ID: ''+signUpRecord.id, NAME: params.name, EMAIL: params.email, DATE: moment().format('ddd, MMM D, YYYY HH:mm Z'), ABOUT: params.about, APP_VERSION: params.appVersion, APP_PLATFORM: params.appPlatform };
    return sendEmail({ to: 'andrey@tarantsov.com', from: 'bot@livereload.com', replyTo: params.email }, NEW_BETA_SIGNUP_EMAIL, info);
    }).done();

Model:

    function recordBetaSignup(ctx, params) {
      return queryRow(ctx, "INSERT INTO beta_signups (name, email, about, app_platform, app_version) VALUES ($1, $2, $3, $4, $5) RETURNING id", params.name, params.email, params.about, params.appPlatform, params.appVersion);
    }


## TODO: Compaign Monitor subscribe hook

    <?php

    $body = http_get_request_body();

    $result = json_decode($body);

    foreach($result->Events as $event) {
      $ip = $event->SignupIPAddress;
      $email = $event->EmailAddress;
      $date  = $event->Date;
      $about = '';
      $notified = FALSE;
      foreach ($event->CustomFields as $field) {
        if ($field->Key == 'About') {
          $about = $field->Value;
        }
      }
      if (!empty($about)) {
        $msg = "Subscriber: $email\nDate: $date\n\nMessage:\n$about\n\n-- LiveReload";
        $notified = mail('andrey@tarantsov.com', "LiveReload subscription: $email", $msg, "From: notification@livereload.com\r\nReply-To: $email");
      }
    }

    ?>


## TODO: compiler messages

    <?php

    require_once("NFSN/RemoteAddr.php");

    // CREATE TABLE unparsable_logs(id int primary key auto_increment, time int not null, date date not null, ip varchar(100) not null, version varchar(100) not null, iversion varchar(100) not null, agent varchar(255) not null, compiler varchar(100) not null, body text not null, index (date))

    $version = $_GET['v'];
    $iversion = $_GET['iv'];
    $compiler = $_GET['compiler'];
    $ip = LastRemoteAddr();
    $agent = urldecode($_SERVER['HTTP_USER_AGENT']);
    $time = time();
    $body = http_get_request_body();

    if (!empty($version) && !empty($iversion) && !empty($compiler)) {
        include '../dbconfig.php';

        $sql = sprintf('INSERT INTO unparsable_logs(time, date, ip, version, iversion, agent, compiler, body) VALUES(%s, FROM_UNIXTIME(%s), "%s", "%s", "%s", "%s", "%s", "%s")',
            $time, $time,
            mysql_real_escape_string($ip),
            mysql_real_escape_string($version),
            mysql_real_escape_string($iversion),
            mysql_real_escape_string($agent),
            mysql_real_escape_string($compiler),
            mysql_real_escape_string($body));

        if (!mysql_query($sql)) {
          header("500 Server Error\r\n");
          die("Internal error, saving failed: " . mysql_error());
        }

        $id = mysql_insert_id();
        $msg = "Unparsable log record $id for compiler $compiler.\n\nIP: $ip\nUser Agent: $agent\n\nLog:\n$body\n\n-- LiveReload";
        mail('andreyvit@me.com', "[LiveReload] Unparsable log for $compiler", $msg, "From: notification@livereload.com\r\nReply-To: andreyvit@me.com");

        die("OK.");
    } else {
      header("400 Bad Request\r\n");
      die("Bad request.");
    }


## TODO: schema

    CREATE TYPE WIDGET_TYPE AS ENUM('bolt', 'screw');

    CREATE TABLE products (
        code VARCHAR PRIMARY KEY,
        name VARCHAR NOT NULL,

        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    INSERT INTO products(code, name) VALUES ('LR', 'LiveReload');

    CREATE TABLE license_types (
        code VARCHAR PRIMARY KEY,
        name VARCHAR NOT NULL,

        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    INSERT INTO license_types(code, name) VALUES ('A', 'Personal License');
    INSERT INTO license_types(code, name) VALUES ('B', 'Business License');
    INSERT INTO license_types(code, name) VALUES ('E', 'Site License');

    CREATE TABLE stores (
        code VARCHAR PRIMARY KEY,
        name VARCHAR NOT NULL,

        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    INSERT INTO stores(code, name) VALUES ('StackSocial', 'StackSocial.com');

    CREATE TABLE licenses (
        id SERIAL PRIMARY KEY,
        product_code VARCHAR NOT NULL REFERENCES products ON DELETE RESTRICT,
        license_type VARCHAR NOT NULL REFERENCES license_types ON DELETE RESTRICT,
        license_code VARCHAR NOT NULL UNIQUE,

        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

        claimed BOOLEAN NOT NULL DEFAULT FALSE,
        claimed_at TIMESTAMPTZ NULL,
        claim_store VARCHAR NULL REFERENCES stores ON DELETE RESTRICT,
        claim_txn VARCHAR NULL,
        claim_qty INTEGER NULL,
        claim_ticket VARCHAR NULL,
        claim_company VARCHAR NULL,
        claim_first_name VARCHAR NULL,
        claim_last_name VARCHAR NULL,
        claim_email VARCHAR NULL,
        claim_notes TEXT NULL
    );

    CREATE TYPE PLATFORM_TYPE AS ENUM('mac', 'windows');

    CREATE TABLE beta_signups (
        id SERIAL PRIMARY KEY,
        received_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

        name VARCHAR NOT NULL,
        email VARCHAR NOT NULL,
        about TEXT NOT NULL,

        app_platform PLATFORM_TYPE NOT NULL,
        app_version VARCHAR NOT NULL
    );
    ALTER TABLE licenses ADD COLUMN claim_full_name VARCHAR NULL;
    ALTER TABLE licenses ADD COLUMN claim_message TEXT NULL;

    INSERT INTO stores(code, name) VALUES ('Paddle', 'Paddle');
    INSERT INTO stores(code, name) VALUES ('manual', '(Manual)');

    ALTER TABLE licenses ADD COLUMN claim_raw JSONB NULL;
    ALTER TABLE licenses ADD COLUMN claim_currency VARCHAR NULL;
    ALTER TABLE licenses ADD COLUMN claim_price NUMERIC(10,4) NULL;
    ALTER TABLE licenses ADD COLUMN claim_sale_gross NUMERIC(10,4) NULL;
    ALTER TABLE licenses ADD COLUMN claim_sale_tax NUMERIC(10,4) NULL;
    ALTER TABLE licenses ADD COLUMN claim_processor_fee NUMERIC(10,4) NULL;
    ALTER TABLE licenses ADD COLUMN claim_earnings NUMERIC(10,4) NULL;
    ALTER TABLE licenses ADD COLUMN claim_coupon VARCHAR NULL;
    ALTER TABLE licenses ADD COLUMN claim_coupon_savings NUMERIC(10,4) NULL;
    ALTER TABLE licenses ADD COLUMN claim_additional VARCHAR NULL;
    ALTER TABLE licenses ADD COLUMN claim_country VARCHAR NULL;


## TODO: Tickets

Example:

    {
      "name": "John Appleseed",
      "email": "andrey@tarantsov.com",
      "product": "livereload-app",
      "platform": "mac",
      "category": "problem",
      "problem": "crash",
      "body": "Hello, world!!!",
      "urgency": "week"
    }

Route:

    var DetailFields = ['subject', 'category', 'problem', 'urgency', 'weblanguage', 'webframework', 'platform', 'product'];
    var DetailFieldLabels = {
      'subject': 'Subject',
      'category': 'Category',
      'problem': 'Problem',
      'urgency': 'Urgency',
      'weblanguage': 'Web-Language',
      'webframework': 'Web-Framework',
      'platform': 'Platform',
      'product': 'Product'
    };

    exports.submit = function(req, res) {
      var details = req.body;
      if (!details.email && req.query.email)
        details = req.query;

      var name = details.name || 'User';
      var email = details.email || '';
      var subject = details.subject || '';
      var product = details.product || '';
      var message = details.body || ''; delete details.body;
      var category = details.category || '';
      var problem = details.problem || '';
      var urgency = details.urgency || '';

      if (!email) {
        res.jsonp({ ok: false, code: 'EINVALID', field: 'email', message: "Missing a required field" });
        return;
      }

      var ticketId = '' + (1000 + Math.floor(Math.random() * 1000));

      var subjectItems = [];
      subjectItems.push('Ticket #' + ticketId);

      if (urgency === 'urgent')
        subjectItems.push('[URGENT]');
      else if (urgency != 'week')
        subjectItems.push('[' + urgency + ']');

      if (product && product !== 'livereload-app')
        subjectItems.push(product + ':');

      if (category === 'problem')
        subjectItems.push('[pr]');
      else if (category === 'feature-request')
        subjectItems.push('[feature]');
      else if (category)
        subjectItems.push('[' + category + ']');

      if (problem && problem !== 'other')
        subjectItems.push('[' + problem + ']');

      subjectItems.push(subject);

      var subject = subjectItems.join(' ');

      var detailItems = [];
      detailItems.push(['From', name + ' <' + email + '>']);
      DetailFields.forEach(function(field) {
        if (details[field]) {
          detailItems.push([DetailFieldLabels[field] || field, details[field]]);
        }
      });
      for (var field in details) if (details.hasOwnProperty(field)) {
        if (field === 'callback' || field === '_')
          continue;  // JSONP and jQuery cache busting
        if (field === 'name' || field === 'email')
          continue;  // handled separately
        if (DetailFields.indexOf(field) === -1) {
          detailItems.push([DetailFieldLabels[field] || field, details[field]]);
        }
      }

      detailItems.push(['User-Agent', req.useragent.source]);
      detailItems.push(['Detected-Browser', req.useragent.Browser + ' ' + req.useragent.Version]);
      detailItems.push(['Detected-OS', req.useragent.OS]);

      var detailRows = [];
      detailItems.forEach(function(item) {
        detailRows.push(item[0] + ": " + item[1] + "\n");
      });

      var body = message + "\n\n---\n\n" + detailRows.join('');

      var userBody = "Your ticket has been emailed to Andrey Tarantsov; here's a copy for your convenience. (You can reply to add more details as well.)\n\n---\n\n" + body;

      postmark.send({
        To: "LiveReload Support <support@livereload.com>",
        From: name + " <support@livereload.com>",
        ReplyTo: name + " <" + email + ">",
        Subject: subject,
        TextBody: body
      }, function(error, success) {
        if (error) {
          res.jsonp({ ok: false, message: error.message });
        } else {
          postmark.send({
            To: name + " <" + email + ">",
            From: name + " <support@livereload.com>",
            ReplyTo: "LiveReload Support <support@livereload.com>",
            Subject: subject,
            TextBody: userBody
          }, function(error, success) {
            if (error) {
              res.jsonp({ ok: true, userOk: false, message: error.message });
            } else {
              res.jsonp({ ok: true, userOk: true });
            }
          });
        }
      })
    };
