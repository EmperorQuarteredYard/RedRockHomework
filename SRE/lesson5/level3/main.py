# -*- coding: utf-8 -*-
"""
重庆邮电大学教务处通知公告爬取脚本
使用Selenium自动化爬取通知公告的标题和内容

使用说明：
1. 安装依赖：pip install selenium webdriver-manager
2. 确保Chrome浏览器已安装
3. 运行脚本：python main.py
"""

import time
import csv
import os
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException, NoSuchElementException
from webdriver_manager.chrome import ChromeDriverManager
from selenium.webdriver.edge.service import Service


def setup_driver(headless=False):
    """配置并返回Chrome浏览器驱动"""
    chrome_options = Options()
    
    # 基础配置
    chrome_options.add_argument('--disable-gpu')
    chrome_options.add_argument('--window-size=1920x1080')
    chrome_options.add_argument('--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36')
    chrome_options.add_argument('--disable-blink-features=AutomationControlled')  # 防止被检测为自动化工具
    chrome_options.add_experimental_option('excludeSwitches', ['enable-automation'])  # 排除自动化开关
    chrome_options.add_experimental_option('useAutomationExtension', False)  # 禁用自动化扩展
    
    # 无头模式
    if headless:
        chrome_options.add_argument('--headless')
    
    # 使用webdriver-manager自动管理Chrome驱动
    service = Service(ChromeDriverManager().install())
    
    # 初始化驱动
    driver = webdriver.Chrome(service=service, options=chrome_options)
    
    # 绕过webdriver检测
    driver.execute_script("Object.defineProperty(navigator, 'webdriver', {get: () => undefined})")
    
    return driver


def navigate_to_notice_page(driver, base_url):
    """导航到通知公告页面"""
    try:
        driver.get(base_url)
        time.sleep(3)  # 等待页面加载
        
        # 尝试点击通知公告链接，需要根据实际页面结构调整选择器
        # 可能的选择器: 
        # - 根据文本: //a[contains(text(), '通知公告')]
        # - 根据类名: .nav-menu a.notice
        # - 根据ID: #notice-link
        
        try:
            # 尝试多种选择器定位通知公告链接
            notice_link_selectors = [
                "//a[contains(text(), '通知公告')]",
                "//a[contains(@href, 'notice')]",
                "//a[contains(@class, 'notice')]",
                ".nav a:contains('通知公告')"
            ]
            
            notice_link = None
            for selector in notice_link_selectors:
                try:
                    if selector.startswith('//'):  # XPath
                        notice_link = WebDriverWait(driver, 10).until(
                            EC.element_to_be_clickable((By.XPATH, selector))
                        )
                    else:  # CSS
                        notice_link = WebDriverWait(driver, 10).until(
                            EC.element_to_be_clickable((By.CSS_SELECTOR, selector))
                        )
                    break
                except:
                    continue
            
            if notice_link:
                notice_link.click()
                time.sleep(3)  # 等待页面跳转
                print("成功导航到通知公告页面")
                return True
            else:
                print("未找到通知公告链接，可能需要手动调整选择器")
                return False
                
        except Exception as e:
            print(f"点击通知公告链接失败: {e}")
            print("将尝试在当前页面查找公告")
            return True  # 继续执行，尝试在当前页面查找公告
            
    except Exception as e:
        print(f"导航到首页失败: {e}")
        return False


def get_notice_links(driver):
    """获取通知公告页面的所有公告链接"""
    notice_links = []
    
    try:
        # 尝试多种选择器定位公告列表
        notice_list_selectors = [
            'div.news-list',
            'ul.notice-list',
            'div.list-container',
            '#notice-container'
        ]
        
        notice_list = None
        for selector in notice_list_selectors:
            try:
                notice_list = WebDriverWait(driver, 10).until(
                    EC.presence_of_element_located((By.CSS_SELECTOR, selector))
                )
                break
            except:
                continue
        
        if not notice_list:
            # 如果没找到列表容器，尝试直接找所有链接
            notice_elements = driver.find_elements(By.CSS_SELECTOR, 'a[href*="notice"]')
        else:
            # 在列表容器内找链接
            notice_elements = notice_list.find_elements(By.CSS_SELECTOR, 'a')
        
        # 去重并筛选有效的公告链接
        seen_links = set()
        for element in notice_elements:
            try:
                link = element.get_attribute('href')
                if link and link not in seen_links:
                    # 筛选包含公告相关关键词的链接
                    if any(keyword in link.lower() for keyword in ['notice', 'gg', 'gonggao', '公告']):
                        notice_links.append(link)
                        seen_links.add(link)
            except Exception as e:
                print(f"获取链接失败: {e}")
        
        # 如果找到的链接太少，尝试另一种方式
        if len(notice_links) < 5:
            print("找到的公告链接较少，尝试其他方式...")
            all_links = driver.find_elements(By.TAG_NAME, 'a')
            for link in all_links:
                try:
                    href = link.get_attribute('href')
                    text = link.text.strip()
                    if href and text and href not in seen_links:
                        notice_links.append(href)
                        seen_links.add(href)
                except:
                    continue
        
        notice_links = list(seen_links)  # 最终去重
        print(f"找到 {len(notice_links)} 条公告链接")
        
    except Exception as e:
        print(f"获取公告链接失败: {e}")
    
    return notice_links


def get_notice_content(driver, url):
    """获取单条公告的标题和内容"""
    try:
        driver.get(url)
        time.sleep(2)  # 等待页面加载
        
        # 获取标题 - 尝试多种选择器
        title = None
        title_selectors = [
            'h1.title',
            'h2.notice-title',
            '#article-title',
            '.news-title'
        ]
        
        for selector in title_selectors:
            try:
                title_element = driver.find_element(By.CSS_SELECTOR, selector)
                title = title_element.text.strip()
                if title:
                    break
            except:
                continue
        
        if not title:
            # 尝试XPath根据文本内容获取标题
            try:
                title_element = driver.find_element(By.XPATH, "//h1 | //h2 | //h3")
                title = title_element.text.strip()
            except:
                title = "无标题"
        
        # 获取内容 - 尝试多种选择器
        content = None
        content_selectors = [
            'div.content',
            'div.article-content',
            'div.news-content',
            '#content'
        ]
        
        for selector in content_selectors:
            try:
                content_element = driver.find_element(By.CSS_SELECTOR, selector)
                content = content_element.text.strip()
                if content:
                    break
            except:
                continue
        
        if not content:
            # 尝试获取页面主体内容
            try:
                main_content = driver.find_element(By.TAG_NAME, 'main')
                content = main_content.text.strip()
            except:
                try:
                    body_content = driver.find_element(By.TAG_NAME, 'body')
                    content = body_content.text.strip()
                except:
                    content = "无内容"
        
        return {
            'title': title,
            'content': content,
            'url': url
        }
        
    except Exception as e:
        print(f"获取公告内容失败 ({url}): {e}")
        return None


def save_to_csv(notices, filename='notices.csv'):
    """将公告数据保存到CSV文件"""
    if not notices:
        print("没有数据可以保存")
        return
    
    fieldnames = ['title', 'content', 'url']
    
    # 确保文件目录存在
    os.makedirs(os.path.dirname(filename) if os.path.dirname(filename) else '.', exist_ok=True)
    
    try:
        with open(filename, 'w', newline='', encoding='utf-8-sig') as csvfile:
            writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
            writer.writeheader()
            for notice in notices:
                writer.writerow(notice)
        
        print(f"数据已成功保存到 {filename}")
        return True
    except Exception as e:
        print(f"保存数据失败: {e}")
        return False


def main():
    """主函数"""
    base_url = 'https://jw.cqupt.edu.cn/'
    output_file = 'notices.csv'
    
    print("初始化浏览器...")
    driver = setup_driver(headless=False)  # 设置为True可开启无头模式
    
    try:
        print(f"访问网站: {base_url}")
        
        # 导航到通知公告页面
        success = navigate_to_notice_page(driver, base_url)
        if not success:
            print("无法导航到通知公告页面，程序终止")
            return
        
        # 获取公告链接
        print("获取公告链接...")
        notice_links = get_notice_links(driver)
        
        if not notice_links:
            print("未找到任何公告链接，可能需要调整选择器")
            return
        
        # 获取公告内容
        print("获取公告内容...")
        notices = []
        
        # 限制爬取数量，避免给服务器造成过大压力
        max_notices = 10
        notice_links = notice_links[:max_notices]
        
        for i, link in enumerate(notice_links, 1):
            print(f"处理第 {i}/{len(notice_links)} 条公告: {link}")
            notice = get_notice_content(driver, link)
            if notice:
                notices.append(notice)
            
            # 随机延迟，避免被反爬
            time.sleep(1 + (i % 3))
        
        # 保存数据
        print("保存数据...")
        save_to_csv(notices, output_file)
        
        print(f"\n爬取完成！共获取 {len(notices)} 条公告")
        
    except Exception as e:
        print(f"程序执行失败: {e}")
    finally:
        driver.quit()
        print("浏览器已关闭")


if __name__ == "__main__":
    main()