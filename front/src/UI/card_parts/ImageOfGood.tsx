import { useEffect, useState } from 'react'

import { Image } from '@chakra-ui/react'

type ImageOfGoodProps = {
  image_url: string
}

const ImageOfGood = ({ image_url }: ImageOfGoodProps) => {
  const [url, setImgUrl] = useState<string>(image_url)

  useEffect(() => {
    const fetchImage = async () => {
      try {
        const response = await fetch(`http://image_service:80/images/${image_url}`)
        const blob = await response.blob()
        const url = URL.createObjectURL(blob)
        setImgUrl(url)
      } catch (error) {
        console.error('Ошибка при загрузке изображения', error)
      }
    }

    fetchImage()
  }, [])

  return (
    <Image
      width="170px"
      height="170px"
      objectFit="cover"
      mt="10px"
      ml="10px"
      mr="10px"
      borderRadius="10px"
      src={url}
    />
  )
}
export default ImageOfGood
