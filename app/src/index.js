import React from 'react';
import ReactDOM from 'react-dom/client';
import Login from './JS/Login';
import Navigation_Bar from './components/Navigation-bar';
import './CSS/index.css';


export default function Main()  
{
    return (    
        <div>
            <Login />
            <Navigation_Bar Display = "block"/>
        </div>
    );
}
const root = ReactDOM.createRoot(document.getElementById('root'));  
root.render(<Main />);

