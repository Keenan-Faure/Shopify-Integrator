import {useEffect} from 'react';
import '../../CSS/page1.css';

function Detailed_product(props)
{
    useEffect(()=> 
    {
        
    }, []);

    return (

        <div className = "details-details">
            <div className = "details-details">
            <div className = "details-image" style = {{backgroundImage: `url(${props.Product_Image})`}}></div>
            <div className = "detailed">
                <div className = "details-title">{props.Product_Title}</div>
                

            </div>
            
        </div>
        </div>
    );
};

Detailed_product.defaultProps = 
{
    Product_Title: 'Product title',
    Product_Code: 'Product code',
    Product_Options: 'Options',
    Product_Category: 'Category',
    Product_Type: 'Type',
    Product_Vendor: 'Vendor',
    Product_Image: '#ccc',
    Product_Price: 'Price'
}
export default Detailed_product;
/*
<div className = "details-image" style = {{backgroundImage: `linear-gradient(to bottom, transparent, white), url(${props.Product_Image})`}}>
*/