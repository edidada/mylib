cmake_minimum_required(VERSION 3.14)
project(concurrency_server CXX)

# 设置 C++ 标准
set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

# 编译选项
add_compile_options(-Wall -Wextra -O2 -pthread)

# 检测操作系统并添加宏定义
if (${CMAKE_SYSTEM_NAME} STREQUAL "Linux")
    add_compile_definitions(LINUX)
elseif (${CMAKE_SYSTEM_NAME} STREQUAL "Darwin")
    add_compile_definitions(MACOS)
endif()

# 可执行文件名称
set(TARGET concurrency_server)

# 源文件和头文件
set(SOURCES
        main.cpp
        base/event_dispatcher.cpp
        base/coroutine.cpp
)

# 头文件用于依赖管理（仅用于 include，不参与编译）
set(HEADERS
        base/server_base.h
        base/event_dispatcher.h
        base/select_dispatcher.h
        base/poll_dispatcher.h
        base/epoll_dispatcher.h
        base/kqueue_dispatcher.h
        single_process_server.h
        multi_process_server.h
        multi_thread_server.h
        process_pool1_server.h
        process_pool2_server.h
        thread_pool_server.h
        leader_follower_server.h
        select_server.h
        poll_server.h
        epoll_server.h
        kqueue_server.h
        reactor_server.h
        coroutine_server.h
)

# 创建可执行文件
add_executable(${TARGET} ${SOURCES} ${HEADERS})

# 安装路径（可选）
install(TARGETS ${TARGET} DESTINATION bin)

# ================== 测试目标 ==================

# 使用自定义命令模拟 Makefile 中的 test-* 目标
function(add_test_target model port)
    add_custom_target(test-${model}
            COMMAND ./${TARGET} ${model} ${port}
            DEPENDS ${TARGET}
            WORKING_DIRECTORY ${CMAKE_CURRENT_BINARY_DIR}
            COMMENT "Running server with model: ${model}"
    )
endfunction()

# 添加多个测试模型
foreach(model
        single_process
        multi_process
        multi_thread
        process_pool1
        process_pool2
        thread_pool
        leader_follower
        select
        poll
        epoll
        kqueue
        reactor
        coroutine
        work_stealing
        actor
        fiber
        producer_consumer
        half_sync_async
        proactor
        pipeline
        hybrid
)
    add_test_target(${model} 8080)
endforeach()

# 添加测试全部的目标
add_custom_target(test-all)
foreach(model
        single_process
        multi_process
        multi_thread
        process_pool1
        process_pool2
        thread_pool
        leader_follower
        select
        poll
        epoll
        kqueue
        reactor
        coroutine
        work_stealing
        actor
        fiber
        producer_consumer
        half_sync_async
        proactor
        pipeline
        hybrid
)
    add_dependencies(test-all test-${model})
endforeach()

# ================== 性能测试目标 ==================
find_program(AB_PROGRAM ab)

if (AB_PROGRAM)
    message(STATUS "Found Apache Bench (ab): ${AB_PROGRAM}")
    add_custom_target(bench
            COMMAND echo "Starting server in background..."
            COMMAND ${CMAKE_CURRENT_BINARY_DIR}/${TARGET} thread_pool 8080 &
            COMMAND sleep 2
            COMMAND echo "Running benchmark..."
            COMMAND ${AB_PROGRAM} -n 1000 -c 10 http://localhost:8080/
            COMMAND echo "Stopping server..."
            COMMAND pkill -f ${TARGET}
            DEPENDS ${TARGET}
            WORKING_DIRECTORY ${CMAKE_CURRENT_BINARY_DIR}
    )
else()
    message(WARNING "Apache Bench 'ab' not found, bench target will not be available.")
endif()

# ================== 帮助目标 ==================
add_custom_target(show-help)
add_custom_command(TARGET show-help POST_BUILD
        COMMAND echo "Available targets:"
        COMMAND echo "  all              - Build the server"
        COMMAND echo "  clean            - Remove built files"
        COMMAND echo "  test-<model>     - Test specific server model"
        COMMAND echo "  test-all         - Run all test models"
        COMMAND echo "  bench            - Run performance benchmark"
        COMMAND echo "  help             - Show this help"
        COMMAND echo ""
        COMMAND echo "Available server models:"
        COMMAND echo "  single_process, multi_process, multi_thread"
        COMMAND echo "  process_pool1, process_pool2, thread_pool"
        COMMAND echo "  leader_follower, select, poll, epoll, kqueue"
        COMMAND echo "  reactor, coroutine"
        COMMAND echo ""
        COMMAND echo "Usage: ./${TARGET} <model> [port]"
        COMMAND echo "Example: ./${TARGET} thread_pool 8080"
)