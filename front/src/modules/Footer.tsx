import { Box, Flex, Text } from '@chakra-ui/react'

const Footer = () => {
  return (
    <Box>
      <Flex
        borderTopLeftRadius="10px"
        borderTopRightRadius="10px"
        height="40px"
        bg="rgba(52, 72, 255, 1)"
        color="white"
        justifyContent="space-around"
        alignItems="center"
      >
        <Flex flex="1" alignItems="center" justifyContent="space-around">
          <Text>+79990058643</Text>
        </Flex>
        <Flex flex="1" alignItems="center" justifyContent="space-around">
          <Text color="white" fontWeight="bold" fontSize="20px">
            Surius
          </Text>
        </Flex>
        <Flex flex="1" alignItems="center" justifyContent="space-around">
          <Text>renevartemiy@yandex.ru</Text>
        </Flex>
      </Flex>
    </Box>
  )
}

export default Footer
