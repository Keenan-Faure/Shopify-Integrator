import '../../../CSS/detailed.css';

function Product_Variants(props)
{
    return (
        <div style={{marginBottom: '5px'}}>
            <table style={{marginBottom: '25px'}}>
                <tbody>
                    <tr>
                        <th style= {{width: '50%'}}>Variant Barcode</th>
                        <th style= {{width: '50%'}}>Variant SKU</th>
                    </tr>
                    <tr>
                        <td style= {{width: '50%'}} className = "barcode">{props.Variant_Barcode}</td>
                        <td style= {{width: '50%'}} className = "sku">{props.Variant_SKU}</td>
                    </tr>
                </tbody>
            </table>
            <table style={{marginBottom: '15px'}}>
                <tbody>
                    <tr>
                        <th>Option 1</th>
                        <th>Option 2</th>
                        <th>Option 3</th>
                    </tr>
                    <tr>
                        <td className = "option1" >{props.Option1}</td>
                        <td className = "option2" >{props.Option2}</td>
                        <td className = "option3" >{props.Option3}</td>
                    </tr>
                </tbody>
            </table>

            <div className = "vr">
                <div className = "Prices" style = {{textAlign: 'center'}}>Variant Prices:</div>
                <br />
                <div className = "price">{props.Price}</div>

                <div className = "Quantities">Quantities</div>
                <br />
                <div className = "quantities">{props.Quantities}</div>

                <div className = "variant-updateDate">{props.Variant_UpdateDate}</div>
            </div>
        </div>
    );
};

Product_Variants.defaultProps = 
{
    Variant_Title: 'N/A',
    Variant_Barcode: 'N/A',
    Variant_ID: 'N/A',
    Variant_SKU: 'N/A',
    Variant_UpdateDate: 'N/A',
    Price: 'N/A',
    Quantities: 'N/A',
    Option1: 'N/A',
    Option2: 'N/A',
    Option3: 'N/A',
    
}
export default Product_Variants;
