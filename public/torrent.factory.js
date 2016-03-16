(function() {
  'use strict';

  angular.module('app')
  .factory('torrentFactory', torrentFactory);

  torrentFactory.$inject = [];

  function torrentFactory() {
    var services = {
      uploadTorrent: uploadTorrent,
    };

    return services;

    function uploadTorrent() {
      
    }
  }

})();