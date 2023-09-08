import {useEffect} from 'react';
import Page1 from '../components/Page1';
import '../CSS/page1.css';

function Products(props)
{
    useEffect(()=> 
    {
        /* Ensures the navbar is set correctly */
        let navigation = document.getElementById("navbar");
        window.onload = function(event)
        {
            navigation.style.left = "30%";
            navigation.style.position = "absolute";
            navigation.style.width = "70%";
            navigation.style.animation = "MoveLeft 1.2s ease";
        }

    }, []);

    return (
        <>
            <Page1 />
        </>
    );
}

export default Products;