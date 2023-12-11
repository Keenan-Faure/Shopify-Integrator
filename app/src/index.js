import React from 'react';
import ReactDOM from 'react-dom/client';
import NavigationBar from './components/Navigation-bar';
import Auto_Slideshow from './components/Auto-slideshow';
import image from './media/Screenshot.png';
import './CSS/index.css';

export default function Main()  
{
    return (    
        <div>
            <NavigationBar Display = "block"/>
        </div>
    );
}
const root = ReactDOM.createRoot(document.getElementById('root'));  
root.render(<Main />);

