import { library } from '@fortawesome/fontawesome-svg-core';
import { fas } from '@fortawesome/free-solid-svg-icons';
import { faTwitter, faFontAwesome } from '@fortawesome/free-brands-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'; // Import the component

import Home from './Home.jsx';

library.add(fas, faTwitter, faFontAwesome);

function App() {
  return (
    <>
      <Home />
    </>
  );
}

export default App;
