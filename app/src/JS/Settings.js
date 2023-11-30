import {useEffect} from 'react';
import Page2 from '../components/Page2';
import Detailed_order from '../components/semi-components/Order/detailed_order';
import '../CSS/page2.css';

/* Must start with a Caps letter */
function Settings()
{
    useEffect(()=> 
    {
        /* Ensures the navbar is set correctly */
        let navigation = document.getElementById("navbar");
        window.onload = function(event)
        {
            navigation.style.left = "0%";
            navigation.style.position = "relative";
            navigation.style.width = "100%";
        }
        
    }, []);

    return (
        <>
            
            <Page2 />
        </>
    );
}

export default Settings;