/*************************************************
 * Utilities
 *************************************************/
/**
 * Defer execution until condition is met.
 * @param {*} method
 * @returns
 */
 let defer = ( method ) => {
  if( loaded ) { method( ) }
  setTimeout( function( ) { defer( method ) }, 50 );
};

/**
 * Request HTML templates.
 * @param {*} elem
 * @returns
 */
const fetch = ( url ) => {
  return $.ajax({
      context:   this,
      dataType : "html",
      async:     false,
      url :      url
    } ).responseText;
}

/**
 * Load HTML templates.
 * @param {*} elem
 * @returns
 */
const template = ( elem ) => {
  let target = $( elem );
  let src    = target.data( `target` );
  target.html( fetch( `templates/${src}.html` ) );
}


/*************************************************
 * App
 *************************************************/
var App = {
  target: '#app',
  assets: Assets, // /inventory/assets.js
  data( ) {
    return {
      plugins: Plugins // TODO - /assets/plugins.js
    }
  },
  methods: {
    camelCase: ( str ) => {
      let results = "";
      const parts = str.split( ' ' );
      parts.forEach( part => {
        // TODO - Use an Array, then .join() w/ a space character.
        results += part.charAt( 0 ).toUpperCase( ) + part.slice( 1 ) + ' ';
      } )
      return results;
    }
  }
}


/*************************************************
 * Main
 *************************************************/
const Main = ( ) => {
  console.clear();
  // Load {{TEMPLATE}} elements.
  const decorate = ( ) => {
    return new Promise( ( resolve, reject ) => {
      $( `.load` ).each( function( ) {
        template( this );
      } )
      //$( `.load` ).forEach( elem => template( elem ) )
      resolve( );
    } );
  };

  const loadVue = ( ) => {
    let vm = Vue.createApp( App ).mount( App.target );
    //vm.$forceUpdate( );  // Force Vue to redraw.
    // Export Globals => TODO - Move this...
    window.parent.Vue = Vue
    window.parent.vm  = vm;
  }


  decorate( ).then( ( ) => {
    console.log("Templates Loaded...");
    loadVue();
  } );
  // Send a 'init' message.  This tells integration tests that we are ready to go
  //cockpit.transport.wait(function() { });
}


document.addEventListener( "DOMContentLoaded", function ( ) {
  Main( );
} );

