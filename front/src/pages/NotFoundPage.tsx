import { Flex, Image } from '@chakra-ui/react'

import page404 from '../images/404.png'

function Page404() {
  return (
    <Flex mt="200px" justifyContent="center" alignItems="center" className="nasianika-div">
      <Image objectFit="cover" width="800px" height="400px" className="page404" src={page404} />
    </Flex>
  )
}

export default Page404
