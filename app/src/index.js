import React from 'react';
import ReactDOM from 'react-dom/client';
import NavigationBar from './components/Navigation-bar';
import './CSS/index.css';

export default function Main()  
{
    return (    
        <div id = "root">
            <NavigationBar Display = "block"/>
        </div>
    );
}
const root = ReactDOM.createRoot(document.getElementById('root'));  
root.render(<Main />);

