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

                <table>
                    <tbody>
                        <tr>
                            <th>Product_Category Futterkiste</th>
                            <th>Maria Product_Code</th>
                            <th>Product_Type</th>
                            <th>Product_Price</th>

                        </tr>
                        <tr>
                            <td><div className = "detailed-category">Product_Category</div></td>
                            <td><div className = "detailed-code">Product_Code</div></td>
                            <td><div className = "detailed-type">Product_Type</div></td>
                            <td><div className = "detailed-price">Product_Price</div></td>
                        </tr>
                    </tbody>
                </table> 
 

                <div className = "detailed-description">
                    Product Description goes here, and it will be a extremelty long piece of text, of course, this can vary
                    but, on average, it could grow to be this large. But we'll see eyy!~
                </div>
            </div>
            <div className = "details-left"></div>
            <div className = "details-right"></div>
            
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