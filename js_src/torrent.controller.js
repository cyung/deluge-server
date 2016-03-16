(function() {
  'use strict';

  angular.module('app')
  .controller('TorrentCtrl', TorrentCtrl);

  TorrentCtrl.$inject = ['$http', 'Upload'];
  var baseUrl = 'http://localhost:3000';

  function TorrentCtrl($http, Upload) {
    var self = this;
    self.key = 'nothisispatrick';

    self.uploadMagnet = function() {
      $http({
        method: 'POST',
        url: baseUrl + '/magnets',
        data: {
          magnet: self.magnet,
        },
        headers: {
          Authorization: self.key,
        },
      })
      .then(function success() {
        self.magnet = "";
      }).catch(function(err) {
        console.log('err =', err);
      });
    }

    self.uploadTorrent = function(files) {
      console.log('Uploading file...');
      console.log('self.key =', self.key);
      for (var i = 0; i < files.length; i++) {
        Upload.upload({
            url: baseUrl + '/torrents',
            data: {torrent: files[i]},
            headers: {Authorization: self.key}
        }).then(function (resp) {
            console.log('Success');
        }, function (resp) {
            console.log('Error status: ' + resp.status);
        }, function (evt) {

        });
      }
    }
  }

})();