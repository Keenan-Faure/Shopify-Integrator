import {useEffect} from 'react';
import Page2 from '../components/Page2';

import '../CSS/page2.css';
import Detailed_customer from '../components/semi-components/Customer/detailed_customer';

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
            
            <Detailed_customer />
        </>
    );
}

export default Settings;