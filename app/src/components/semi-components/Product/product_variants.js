import '../../../CSS/detailed.css';

function Product_Variants(props)
{
    return (
        <div>
            <div className = "variant-title">{props.Variant_Title}</div>
            <table>
                <tbody>
                    <tr>
                        <th>Variant Barcode</th>
                        <th>Variant ID</th>
                        <th>Variant SKU</th>
                    </tr>
                    <tr>
                        <td>{props.Variant_Barcode}</td>
                        <td>{props.Variant_ID}</td>
                        <td>{props.Variant_SKU}</td>
                    </tr>
                </tbody>
            </table>
            <table>
                <tbody>
                    <tr>
                        <th>Option 1</th>
                        <th>Option 2</th>
                        <th>Option 3</th>
                    </tr>
                    <tr>
                        <td>{props.Option1}</td>
                        <td>{props.Option2}</td>
                        <td>{props.Option3}</td>
                    </tr>
                </tbody>
            </table>

            <div className = "vr">
                <div className = "variant-updateDate">{props.Variant_UpdateDate}</div>
                <div className = "variant-price"><div className = "c">--</div>{props.Price_high}</div><div className = "variant-price">{props.Price_low}</div>
            </div>
        </div>
    );
};

Product_Variants.defaultProps = 
{
    Variant_Title: 'Variant Title',
    Variant_Barcode: 'Variant Barcode',
    Variant_ID: 'Variant ID',
    Variant_SKU: 'Variant SKU',
    Variant_UpdateDate: 'Variant Update Date',
    Price_high: 'High Price',
    Price_low: 'Low Price'
}
export default Product_Variants;
