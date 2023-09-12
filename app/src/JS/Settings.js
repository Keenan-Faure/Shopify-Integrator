import {useEffect} from 'react';
import Page2 from '../components/Page2';
import '../CSS/page2.css';

/* Must start with a Caps letter */
function Settings(props)
{
    useEffect(()=> 
    {
        /* Ensures the navbar is set correctly */
        let navigation = document.getElementById("navbar");
        let model = document.getElementById("model");
        window.onload = function(event)
        {
            navigation.style.left = "0%";
            navigation.style.position = "relative";
            navigation.style.width = "100%";
            model.style.animation = "none";
        }
        
    }, []);

    return (
        <>
            <Page2
            />
            
        </>
    );
}

export default Settings;