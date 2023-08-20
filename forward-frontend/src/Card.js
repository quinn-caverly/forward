import "./App.scss"
import myImage from "./sample.webp"
import { useQuery, gql } from '@apollo/client';

const GET_PRODUCT_BY_ID = gql`
  query GetProductById($Id: String!) {
    getProductById(Id: $Id) {
        Up {
            Brand
            Name
            Id
            UrlOnBrandSite
            ClothingType
        }
        ColorContainers {
            ColorName
            ImageBytes
        }
    }
  }
`;


function buildImage(string) {
    // Convert Base64 string to Uint8Array
    const byteCharacters = atob(string);
    const byteNumbers = new Array(byteCharacters.length);
    for (let i = 0; i < byteCharacters.length; i++) {
        byteNumbers[i] = byteCharacters.charCodeAt(i);
    }
    const uint8Array = new Uint8Array(byteNumbers);

    // Create Blob from Uint8Array
    const blob = new Blob([uint8Array], { type: 'image/jpeg' });

    // Create a URL for the Blob
    const imageUrl = URL.createObjectURL(blob);

    return imageUrl
}

function Card(props) {

  const { loading, error, data } = useQuery(GET_PRODUCT_BY_ID, {
    variables: { Id: props.id },
  });

  if (loading) return (
    <div class="card-holder">
      <div class="card">
      </div>
    </div>
  );

  if (error) return <p>Error: {error.message}</p>;

  const firstColor = data.getProductById.ColorContainers[0]

  const imageUrl = buildImage(firstColor.ImageBytes[0])

  return (
    <div class="card-holder">
      <a href={data.getProductById.Up.UrlOnBrandSite} target="_blank">
      <div class="card">
       <div class="image-holder">
         <img src={imageUrl} class="image"></img>
        </div>
        <div class="brand-name">
          {data.getProductById.Up.Brand}
        </div>
        <div class="text-holder">
            <div class="item-name">
                {data.getProductById.Up.Name}
            </div>
       </div>
      </div>
      </a>
    </div>
  );
}

export default Card;
