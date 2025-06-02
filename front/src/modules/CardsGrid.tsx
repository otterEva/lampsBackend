import { useEffect, useState } from 'react'

import { Grid } from '@chakra-ui/react'

import SingleCard from '../components/SingleCard'

type GoodsData = {
  cost: string
  name: string
  description: string
  image_url: string
}

const CardsGrid = () => {
  const [goodsInfo, setGoodsData] = useState<GoodsData[]>([])
  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch('http://goods_service:80/goods')
        const data = await response.json()
        setGoodsData(data)
      } catch (error) {
        console.error('Ошибка при получении данных:', error)
      }
    }

    fetchData()
  }, [])

  return (
    <Grid
      gridTemplateColumns={{
        base: 'repeat(1, 1fr)',
        sm: 'repeat(2, 1fr)',
        md: 'repeat(3, 1fr)',
        lg: 'repeat(4, 1fr)',
        xl: 'repeat(6, 1fr)',
      }}
      columnGap="20px"
      rowGap="20px"
      justifyItems="center"
      alignItems="top"
    >
      {goodsInfo.map((Good) => (
        <SingleCard
          price={Good.cost}
          key={Good.image_url}
          image_url={Good.image_url}
          title={Good.name}
        />
      ))}
    </Grid>
  )
}

export default CardsGrid
